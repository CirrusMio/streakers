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

### Without Vagrant

Make sure your `GOPATH` is set. Perhaps [just use one GOPATH](http://mwholt.blogspot.com/2014/02/why-i-use-just-one-gopath.html)

Install dependencies:

`go get github.com/nitrous-io/goop`
`goop install`

then

`goop exec go run server.go`

load up [localhost:3000](http://localhost:3000)

If you have your path setup right you can use one of these live reload packages:

Use [Gin](https://github.com/codegangsta/gin) and [Fresh](https://github.com/pilu/fresh)
