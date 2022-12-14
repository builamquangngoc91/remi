FROM golang:latest

RUN mkdir -p /go/src/remi

WORKDIR /go/src/remi

COPY . /go/src/remi

RUN go install remi

CMD /go/bin/remi

EXPOSE 8080