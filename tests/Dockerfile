FROM ubuntu:16.04

RUN apt-get -y -qq update && apt-get -qq -y install \
    python3-pip \
    python3-pytest \
    docker.io \
    docker-compose

RUN pip3 install requests minio

RUN mkdir -p /tests
ENTRYPOINT ["bash", "/tests/run.sh"]
