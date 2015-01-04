FROM stackbrew/ubuntu:trusty

MAINTAINER Nick Warner <nickwarner@gmail.com>

RUN apt-get update
RUN apt-get install -y golang git

ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

RUN go get github.com/nitrous-io/goop
RUN goop install

EXPOSE 3000

ADD * /app/
CMD ["goop exec", "go", "run", "/app/server.go"]
