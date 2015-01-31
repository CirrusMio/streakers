streakers
=========

don't break the chain - commit to open source everyday


## Goals

- Track a user's commit streak
- Remind users if time's running out on their streak
- Watch all major public git hosting services


## Stretch goals

- API
- containerized deployment
- suggest projects to work on
- recommend what to work on next based on favorites
- merit badges
- ask for help on your open source project

## Development

### Developing With Docker

Install and setup Docker on your preferred platform.

If you're on Mac OS X be aware that you need to use the IP address given to you by boot2docker (`echo $DOCKER_HOST`) rather than `localhost` or `127.0.0.1`.

```sh
docker build . -t CirrusMio/streakers

docker run -it --rm -p 8080:8080 --name streakers CirrusMio/streakers
```

### Local Development

Make sure you have Go installed and your `GOPATH` is set.
Read [Getting Started](http://golang.org/doc/install) and [How to Write Go Code](http://golang.org/doc/code.html) to get started.

```sh
cd $GOPATH/src/github.com/CirrusMio/streakers

go get github.com/CirrusMio/streakers
go install github.com/CirrusMio/streakers

$GOPATH/bin/streakers
```

## Resources

[Go Bootcamp Book](http://www.golangbootcamp.com/book)

[Programming in Go](http://www.golang-book.com/)

[Resources for new Go programmers](http://dave.cheney.net/resources-for-new-go-programmers)

[Go Project Structure For Rubyists](http://gofullstack.com/articles/go-project-structure-for-rubyists.html)

[Effective Go](http://golang.org/doc/effective_go.html)

[Learn X in Y minutes Where X=Go](http://learnxinyminutes.com/docs/go/)

[Go by Example](https://gobyexample.com/)

[Gophercasts](https://gophercasts.io)

[Go Programming Language Specification](http://golang.org/ref/spec)

[Go Memory Model](http://golang.org/ref/mem)
