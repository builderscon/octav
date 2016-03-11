# octav

octav - The main API server

[![Build Status](https://travis-ci.org/builderscon/octav/octav.svg?branch=master)](https://travis-ci.org/builderscon/octav/octav)

[![GoDoc](https://godoc.org/github.com/builderscon/octav/octav?status.svg)](https://godoc.org/github.com/builderscon/octav/octav)

# Description

octav is the API server that does the basic CRUD operations on data.
No authentication is done on this component.

# Code generation directions

## If you edited spec/v1/api.json

```
hsup -s /path/to/octav/spec/v1/api.json -d /path/to/octav/octav
```

This will regenerate 

* octav/client/client.go
* octav/validator/validator.go
* octav/octav_hsup.go

If you add more endpoints, you need to write additional `doXXXX` handlers in handlers.go

## If you edited transport structs

By "transport" structs, we mean structs like `LookupVenueRequest`, which are used to transport data between the client and the API server.

First, see `How To Build The Code Generation Tools`, and build the tools.
Then run:

```
./gentransport ... -t NewTransportType -d .
```

You have to specify all of the necessary transport types. See `octav.go`'s `go:generate` lines.



## If you edited model types

Structs like `Conference`, `Venue`, `Room`, etc are models. They are the frontends to the underlying DB structs.

First, see `How To Build The Code Generation Tools`, and build the tools.
Then run:

```
./genmodel -t Room -t User -t Venue -d .
```

Note that we don't generate `Conference` here, because `Conference` has some special case handling that can't be automatically generated. We *may* fix this later, but for now, you will have to do this by hand.

You have to specify all of the necessary model types. See `octav.go`'s `go:generate` lines.

## If you edited DB types

Structs defined in `db/interface.go` talk directly to the database.

First, see `How To Build The Code Generation Tools`, and build the tools.
Then run:

```
./gendb -t Conference -t Room -t Session -t User -t Venue -t LocalizedString -d db
```

They generate basic database access routines.

You have to specify all of the necessary db types. See `octav.go`'s `go:generate` lines.

## Running all of the generation tools

You can run `go generate` to generate the DB tools, Model tools, and Transport tools.

## How To Build The Code Generation Tools

```
go build ./internal/cmd/gendb
go build ./internal/cmd/genmodel
go build ./internal/cmd/gentransport
```

## Useful Debugging Tips

### Enable Debug Prints

When running your tests, use the tag `debug0` (or debug). See [github.com/lestrrat/go-pdebug](https://github.com/lestrrat/go-pdebug) for details. |

## Useful Environment Variables

| Name | Description |
|:-----|:------------|
|OCTAV_TEST_DSN | DSN to use to connect to the database. See [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) for syntax details. |
|OCTAV_DB_TRACE | Enable to see SQL execution trace. See [github.com/shogo82148/go-sql-proxy](https://github.com/shogo82148/go-sql-proxy) for details. |
|PDEBUG_TRACE | Enable print debug mode. See [github.com/lestrrat/go-pdebug](https://github.com/lestrrat/go-pdebug) for details. |


# Running Tests

See "Useful Environemnt Variables", and "Useful Debugging Tips"

## 0. Drop the old database

The schema is still changing wildly. You probably want to flush it from
time to time

```
mysqladmin -uroot drop octav
```

## 1. Create a database

```
mysqladmin -uroot create octav
mysql -uroot octav < octav/sql/octav.sql
```

## 2. Run

```
cd octav
OCTAV_TEST_DSN='root:@:/octav?parseTime=true' go test .
```

Currently, the tests are still failing.

