// 8.12-jessie
FROM golang:latest 

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
CMD ["/app/main"]


------------------------

FROM golang:1.12-alpine

ENV DEP_URL="https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64"
ENV GOPATH=/go

COPY src/   /go/src/github.com/instana/dispatch

WORKDIR /go/src/github.com/instana/dispatch

RUN curl -fsSL -o /usr/local/bin/dep ${DEP_URL} && \
    chmod +x /usr/local/bin/dep && \
    dep ensure && \
    go build -o bin/gorcv

# @todo stage this build

EXPOSE 5000

CMD bin/gorcv

## $ docker build -t go-docker-optimized -f Dockerfile.multistage .
-------------------