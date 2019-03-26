FROM golang:1.11 as builder
ENV GO111MODULE=on 
RUN mkdir -p /go/src/github.com/mendersoftware/mender-cli
WORKDIR /go/src/github.com/mendersoftware/mender-cli
ADD ./ .
RUN go mod download
RUN make build

FROM busybox
COPY --from=builder /go/src/github.com/mendersoftware/mender-cli/mender-cli /
ENTRYPOINT ["/mender-cli"]