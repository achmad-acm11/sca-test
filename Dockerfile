FROM golang:1.21.12-alpine3.20 AS builder
RUN go env -w GO111MODULE=on
WORKDIR /sca-integrator
COPY ./    ./
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.20 as trivy-downloader
WORKDIR /trivy-binary
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
RUN apk add git
RUN apk add curl
RUN curl -LJO https://github.com/aquasecurity/trivy/releases/download/v0.56.2/trivy_0.56.2_Linux-64bit.tar.gz
RUN tar xzvf trivy_0.56.2_Linux-64bit.tar.gz

FROM alpine:3.20 as golang-app
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
COPY --from=trivy-downloader /trivy-binary/trivy ./usr/local/bin/trivy
RUN chmod 755 /usr/local/bin/trivy
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
WORKDIR /root/
COPY --from=builder /sca-integrator ./
#COPY --from=builder /sca-integrator/_public_key.pem ./
RUN mkdir "_project-misconfig-file"
RUN mkdir "_project-repository"
CMD ["./main"]
