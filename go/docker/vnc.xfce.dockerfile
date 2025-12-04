FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive \
    VNC_DISPLAY=:1 \
    VNC_GEOMETRY=1600x900 \
    VNC_DEPTH=24

RUN apt-get update && apt-get install -y --no-install-recommends \
      tigervnc-standalone-server \
      tigervnc-tools \
      tini \
      xfce4 xfce4-terminal dbus-x11 \
      firefox-esr fonts-dejavu-core

RUN apt-get install -y \
    sudo \
    && rm -rf /var/lib/apt/lists/*

COPY scripts/xstartup /xstartup
COPY scripts/entrypoint.sh /entrypoint.sh
COPY scripts/vncserver.sh /vncserver.sh

ENTRYPOINT ["/usr/bin/tini", "-g", "--"]
CMD ["/entrypoint.sh"]
