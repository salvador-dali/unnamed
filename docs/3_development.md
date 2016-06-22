[Setup Go with Pycharm](http://stackoverflow.com/a/37698196/1090562)

### Selection of a routing library

Requirements:
 
 - GET, POST, DELETE, PUT verbs
 - complex routes: `/user/id/animal/id`
 - as lightweight as possible
 - not a framework. Only router

Tried a couple of routers:
 
 - standard Go router is super slow, does not support complex routes
 - [HttpRouter](https://github.com/julienschmidt/httprouter) - does not support 
 [normal routes](https://github.com/julienschmidt/httprouter/issues/12)
 - [Denco](https://github.com/naoina/denco) - does not support DELETE, PUT

The one that supports the route structure and all HTTP Verbs is 
[httptreemux](https://github.com/dimfeld/httptreemux). It is also one of the [fastest and with 
small memory allocation](https://github.com/dimfeld/go-http-routing-benchmark).

### Style conventions

 - run `go fmt`, `go vet`, `golint` before each commit.
 - no underscores in variables, only camelCase
 - comments before every function. Do start with: 'This function analyses ...'. Just 'Analyses ...' 
 - comments inside function should explain why something is done

###  Tests

To run a test, run `go test ./folder` or go to that directory and run `go test`.
To run a single test, run `go test -run TestName`. You can add `-v` to see more details. 

### How to write tests

Write tests only if you see that you have to test manually to often or you are afraid to break something.

Use [table-driven](https://github.com/golang/go/wiki/TableDrivenTests) tests. Sometimes it makes 
sense to create two separate tables: `tableSuccess` and `tableFail`. Do not run test from maps, only
from slices (maps do not run deterministically)

    tableSuccess := []struct {
		field1 type
		field2 type
		...
		fieldN type
	}{
		{3, ..., 2},
		...,
		{1, ..., 1},
	}
	for _, v := range tableSuccess {
	    ...
	}
