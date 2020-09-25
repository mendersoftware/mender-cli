FROM golang:1.15-alpine3.12 as builder
RUN apk update && apk add make git
RUN mkdir -p /go/src/github.com/mendersoftware/mender-cli
WORKDIR /go/src/github.com/mendersoftware/mender-cli
ADD ./ .
RUN make build

FROM busybox
COPY --from=builder /go/src/github.com/mendersoftware/mender-cli/mender-cli /
ENTRYPOINT ["/mender-cli"]
