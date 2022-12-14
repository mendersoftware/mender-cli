FROM tiangolo/docker-with-compose

RUN apk add bash python3 libc6-compat xz-dev
RUN pip3 install requests minio pytest

RUN mkdir -p /tests
ENTRYPOINT ["bash", "/tests/run.sh"]
