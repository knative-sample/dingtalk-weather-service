# Build the manager binary
FROM registry.cn-hangzhou.aliyuncs.com/knative-sample/golang:1.12.9 as builder

# Copy in the go src
WORKDIR /go/src/github.com/knative-sample/dingtalk-weather-service
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o dingtalk-weather-service-receive github.com/knative-sample/dingtalk-weather-service/cmd/receive

# Copy the dingtalk-weather-service-receive into a thin image
FROM registry.cn-beijing.aliyuncs.com/knative-sample/centos:7.6.1810
WORKDIR /
COPY --from=builder /go/src/github.com/knative-sample/dingtalk-weather-service/dingtalk-weather-service-receive app/
ENTRYPOINT ["/app/dingtalk-weather-service-receive"]