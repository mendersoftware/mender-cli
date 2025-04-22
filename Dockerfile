FROM golang:1.24.2 as builder
ARG BUILDFLAGS=
WORKDIR /build
RUN --mount=type=bind,source=.,dst=/build \
  make build BUILDFLAGS="-o /mender-cli ${BUILDFLAGS}"

FROM busybox:1.36.1
COPY --from=builder /mender-cli /usr/bin/
ENTRYPOINT ["/usr/bin/mender-cli"]
