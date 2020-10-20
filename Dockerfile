FROM golang:1.15.3-alpine3.12 as builder
RUN apk add --no-cache make git
WORKDIR /go/src/github.com/mendersoftware/mender-cli
ADD ./ .
RUN make build

FROM busybox
COPY --from=builder /go/src/github.com/mendersoftware/mender-cli/mender-cli /
ENTRYPOINT ["/mender-cli"]
