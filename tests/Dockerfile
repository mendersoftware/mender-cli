# QA-507: Pin to Alpine 3.16 to have OpenSSL 1.1 compatibility while we redesign this dependency
FROM tiangolo/docker-with-compose:2022-11-28

RUN apk add bash python3 libc6-compat xz-dev
RUN pip3 install requests minio pytest

RUN mkdir -p /tests
ENTRYPOINT ["bash", "/tests/run.sh"]
