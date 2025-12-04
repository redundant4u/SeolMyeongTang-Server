import os
import time
from typing import Optional

from kubernetes import config
from kubernetes.client import CoreV1Api

# ── 환경변수(필요에 따라 조정) ──────────────────────────────────────────────
NAMESPACE = os.getenv("K8S_NAMESPACE", "vnc")
LABEL_APP_KEY = os.getenv("LABEL_APP_KEY", "app")
LABEL_APP_VAL = os.getenv("LABEL_APP_VAL", "vnc")
LABEL_SESS_KEY = os.getenv("LABEL_SESSION_KEY", "session-id")
TARGET_PORT = os.getenv("TARGET_PORT", "5901")
CACHE_TTL = int(os.getenv("CACHE_TTL", "10"))
# ────────────────────────────────────────────────────────────────────────────


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


def _select_pod(pods) -> Optional[str]:
    for p in pods:
        if _is_pod_ready(p):
            return p.status.pod_ip

    for p in pods:
        if p.status.phase == "Running" and p.status.pod_ip:
            return p.status.pod_ip
    return None


def _find_target_hostport(token: str) -> Optional[str]:
    now = time.time()
    ent = _cache.get(token)
    if ent and ent[0] > now:
        return ent[1]

    selector = f"{LABEL_APP_KEY}={LABEL_APP_VAL},{LABEL_SESS_KEY}={token}"

    try:
        pods = _core.list_namespaced_pod(
            namespace=NAMESPACE,
            label_selector=selector,
            _request_timeout=3,
        ).items
    except Exception as e:
        print(e)
        return None

    ip = _select_pod(pods)
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
