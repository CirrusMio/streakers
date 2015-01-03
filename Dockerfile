FROM stackbrew/ubuntu:trusty

MAINTAINER Nick Warner <nickwarner@gmail.com>

RUN apt-get update

RUN apt-get install -y golang
RUN apt-get install -y git

ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

RUN go get github.com/codegangsta/negroni
RUN go get github.com/gorilla/mux

EXPOSE 3000

ADD * /app/
CMD ["go", "run", "/app/server.go"]
