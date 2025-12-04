#!/bin/sh

set -e

: "${USER_NAME:=app}"
: "${USER_PASSWORD:=app}"

useradd -m -s /bin/bash ${USER_NAME}

usermod -aG sudo ${USER_NAME}
echo "${USER_NAME}:${USER_PASSWORD}" | chpasswd
unset USER_PASSWORD

exec runuser -u ${USER_NAME} -- /vncserver.sh
