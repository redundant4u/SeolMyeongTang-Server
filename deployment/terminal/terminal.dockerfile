FROM ubuntu:22.04

RUN useradd -U -u 999 -d /home/seol -m -l seol && \
    mkdir -p /home/seol/bin && \
    cp /bin/ls /home/seol/bin/ls && \
    cp /bin/cat /home/seol/bin/cat && \
    cp /bin/clear /home/seol/bin/clear && \
    rm /home/seol/.bash_logout /home/seol/.profile

COPY bashrc /home/seol/.bashrc
COPY about /home/seol/about

USER seol

WORKDIR /home/seol
