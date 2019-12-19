[![Build Status](https://travis-ci.com/kernle32dll/emissione-go.svg?branch=master)](https://travis-ci.com/kernle32dll/emissione-go)
[![GoDoc](https://godoc.org/github.com/kernle32dll/emissione-go?status.svg)](http://godoc.org/github.com/kernle32dll/emissione-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kernle32dll/emissione-go)](https://goreportcard.com/report/github.com/kernle32dll/emissione-go)
[![codecov](https://codecov.io/gh/kernle32dll/emissione-go/branch/master/graph/badge.svg)](https://codecov.io/gh/kernle32dll/emissione-go)

# emissione-go

emissione-go is a small (no dependencies!) framework, which provide dynamic switching support for http response content types.

E.g. this allows you to transparently serve both XML and JSON output in your API.

All this is controlled via the [Accept](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept) header.

emissione-go is inspired by [senlinms/gores](https://github.com/senlinms/gores) and the render types of [gin-gonic/gin](https://github.com/gin-gonic/gin).

Download:

```
go get github.com/kernle32dll/emissione-go
```

Detailed documentation can be found on [GoDoc](https://godoc.org/github.com/kernle32dll/emissione-go).

## Getting started

emissione-go provides two ways for getting started:

You can either use a opinionated default handler by using `emissione.Default()`, or define one yourself via `emissione.New(...)`.
The latter allows you to define a custom mapping, and a default handler. You can look at the source of `emissione.Default()` to get an idea.

Here is a quick example, using the defaults:

```go
package main

import (
	"github.com/kernle32dll/emissione-go"

	"log"
	"net/http"
)

// User is a just sample struct for showcasing.
type User struct {
	Name string `json:"name",xml:"Name"`
}

func main() {
	router := http.NewServeMux()

	em := emissione.Default()

	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		em.Write(w, r, http.StatusOK, User{Name: "Björn Gerdau"})
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
```
Use the following curl calls, to see the code in action:

`curl localhost:8080/user`

`curl -H "Accept: application/xml" localhost:8080/user`

## Disabling stream support

As of 1.1.0, the emissione simple writer (and its implementations) has gained support for streaming responses.

It is now possible to configure a simple writer to use streaming via `emissione.StreamMethod(...)`.
If configured, the streaming method will take precedence over the "old" marshall method.

Per default, all provided writer implementations, and thus the emissione default handler, too, are now configured
**with** streaming. It can be disabled by using `emissione.StreamMethod(nil)`, anywhere an option argument is
expected. For example, `emissione.Default(emissione.StreamMethod(nil))` creates the default handler with
streaming disabled.

For implementation details, it is worthwhile to look at the [JSON writer](./writer_json.go) implementation.

## Extending emissione

Extending emissione is straight forward. Simple use `emissione.New(...)` to define a custom mapping, and if necessary implement
the emissione `Writer` interface. A simple implementation used by emissione itself is `SimpleWriter`, which simply implements
the `Writer` interface by delegating to a marshalling method, and setting the appropriate `Content-Type` header. This allows
for drop-in usage of Go's own marshall methods, such as `json.Marshal` and `xml.Marshal`, or [jsoniter](https://github.com/json-iterator/go).