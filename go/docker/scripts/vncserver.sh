#!/bin/sh

set -e

: "${VNC_DISPLAY:=:1}"
: "${VNC_GEOMETRY:=1600x900}"
: "${VNC_DEPTH:=24}"

mkdir -p "${HOME}/.vnc"

vncserver -kill "${VNC_DISPLAY}" >/dev/null 2>&1 || true

exec vncserver "${VNC_DISPLAY}" \
  -fg \
  -localhost no \
  -geometry "${VNC_GEOMETRY}" \
  -depth "${VNC_DEPTH}" \
  -SecurityTypes None \
  --I-KNOW-THIS-IS-INSECURE
