FROM golang:alpine as builder
WORKDIR /weave
COPY ["main.go", "go.mod", "go.sum", "./"]
COPY docs/ docs/
COPY pkg/ pkg/
COPY static/ static/
RUN go build

FROM alpine
COPY --from=builder /weave/weave /
ENTRYPOINT ["/weave"]
