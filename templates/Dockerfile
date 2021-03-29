ARG GO_VERSION=1.16.2-alpine3.13
FROM golang:$GO_VERSION as builder
WORKDIR /go/src/git.chotot.org/fse/orchestration-svc

COPY go.mod go.sum ./
RUN go mod download
COPY . /go/src/git.chotot.org/fse/orchestration-svc

EXPOSE 8080

RUN go install /go/src/git.chotot.org/fse/orchestration-svc

COPY ./run_all.sh /usr/bin/
CMD ["sh","/usr/bin/run_all.sh"]