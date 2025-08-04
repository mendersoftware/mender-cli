FROM golang:1.24.5 as builder
WORKDIR /build
ARG BUILDFLAGS=
RUN --mount=type=bind,source=.,dst=/build,ro \
  GIT_DISCOVERY_ACROSS_FILESYSTEM=1 make build "BUILDFLAGS=-o /mender-cli ${BUILDFLAGS}"

FROM python:3.13-slim as acceptance

RUN pip3 install requests pytest cryptography boto3 docker
COPY --from=builder /mender-cli /usr/bin/

WORKDIR /tests
ENTRYPOINT ["python3", "-m", "pytest"]
CMD ["tests"]

FROM busybox:1.37.0
COPY --from=builder /mender-cli /usr/bin/
ENTRYPOINT ["/usr/bin/mender-cli"]
