FROM tiangolo/docker-with-compose

RUN apk add bash python3 openssl1.1-compat gcompat
RUN pip3 install requests minio pytest

RUN mkdir -p /tests
ENTRYPOINT ["bash", "/tests/run.sh"]
