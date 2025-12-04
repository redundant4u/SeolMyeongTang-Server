FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive \
    DISPLAY=:1 \
    GEOMETRY=1600x900 \
    DEPTH=24

RUN apt-get update && apt-get install -y --no-install-recommends \
      tigervnc-standalone-server \
      tigervnc-tools \
      tini \
      lxde-core lxterminal dbus-x11 \
    && rm -rf /var/lib/apt/lists/*

COPY scripts/xstartup /home/app/.vnc/xstartup
COPY scripts/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/usr/bin/tini", "-g", "--"]
CMD ["/entrypoint.sh"]
