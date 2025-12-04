import os
import time
from typing import Optional

from kubernetes import config
from kubernetes.client import CoreV1Api

NAMESPACE = os.getenv("K8S_NAMESPACE")
LABEL_APP_KEY = os.getenv("LABEL_APP_KEY")
LABEL_APP_VAL = os.getenv("LABEL_APP_VAL")
LABEL_CLIENT_ID = os.getenv("LABEL_CLIENT_ID")
TARGET_PORT = os.getenv("TARGET_PORT")
CACHE_TTL = int(os.getenv("CACHE_TTL", "10"))


def _load_kube():
    try:
        config.load_incluster_config()
    except Exception:
        config.load_kube_config()
    return CoreV1Api()


_core: CoreV1Api = _load_kube()
_cache = {}


def _is_pod_ready(pod) -> bool:
    if pod.status.phase != "Running" or not pod.status.pod_ip:
        return False

    conds = pod.status.conditions or []
    for c in conds:
        if getattr(c, "type", "") == "Ready" and getattr(c, "status", "") == "True":
            return True
    return False


def _select_pod(pod) -> Optional[str]:
    if _is_pod_ready(pod):
        return pod.status.pod_ip

    return None


def _find_target_hostport(token: str) -> Optional[str]:
    now = time.time()
    ent = _cache.get(token)
    if ent and ent[0] > now:
        return ent[1]

    ids = token.split(":", 1)
    client_id, session_id = ids[0], ids[1]

    selector = f"{LABEL_APP_KEY}={LABEL_APP_VAL},{LABEL_CLIENT_ID}={client_id}"

    try:
        pods = _core.list_namespaced_pod(
            namespace=NAMESPACE,
            label_selector=selector,
            _request_timeout=3,
        ).items
        filtered_pod = next((p for p in pods if p.metadata.name == session_id), None)
    except Exception as e:
        print(e)
        return None

    ip = _select_pod(filtered_pod)
    if not ip:
        return None

    hostport = [ip, TARGET_PORT]
    _cache[token] = (now + CACHE_TTL, hostport)

    return hostport


def lookup(token: str) -> Optional[str]:
    return _find_target_hostport(token)


class TokenPlugin:
    def __init__(self, src: str | None = None):
        self.src = src

    def lookup(self, token: str) -> Optional[str]:
        return _find_target_hostport(token)
