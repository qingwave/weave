FROM golang:alpine  as builder
RUN apk update && \
    apk add git gcc
ADD ./ /go/src/v9.git.n.xiaomi.com/ContainerCloud/k8s-pod-prober/
RUN cd /go/src/v9.git.n.xiaomi.com/ContainerCloud/k8s-pod-prober && \
    export GO111MODULE=on && \
    CGO_ENABLED=0 go install -v v9.git.n.xiaomi.com/ContainerCloud/k8s-pod-prober/cmd/k8s-pod-prober

FROM alpine
COPY --from=builder /go/bin /go/bin
EXPOSE 9798
ENTRYPOINT ["/go/bin/k8s-pod-prober"]