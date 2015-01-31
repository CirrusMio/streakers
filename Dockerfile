FROM golang

MAINTAINER Nick Warner <nickwarner@gmail.com>

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/CirrusMio/streakers

# Change directory so static assets will be in the correct location.
WORKDIR /go/src/github.com/CirrusMio/streakers

# Set up the dependencies for streakers.
RUN go get github.com/CirrusMio/streakers

# Build streakers.
RUN go install github.com/CirrusMio/streakers

# Run streakers.
ENTRYPOINT ["/go/bin/streakers"]

# Expose streakers on 8080.
EXPOSE 8080
