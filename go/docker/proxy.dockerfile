FROM python:3.14-slim

WORKDIR /root

COPY proxy .
RUN pip3 install setuptools numpy kubernetes && \
    python3 setup.py install

ENTRYPOINT ["./run"]
CMD ["0.0.0.0:8080", "--token-plugin", "k8s_plugin.TokenPlugin", "--heartbeat", "30"]
