FROM golang:1.16.5-alpine3.12 as builder
RUN apk add --no-cache make git
WORKDIR /go/src/github.com/mendersoftware/mender-cli
ADD ./ .
RUN make build

FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM busybox
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/github.com/mendersoftware/mender-cli/mender-cli /
ENTRYPOINT ["/mender-cli"]
