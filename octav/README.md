# octav

octav - The main API server

[![Build Status](https://travis-ci.org/builderscon/octav.svg?branch=master)](https://travis-ci.org/builderscon/octav)

[![GoDoc](https://godoc.org/github.com/builderscon/octav/octav?status.svg)](https://godoc.org/github.com/builderscon/octav/octav)

# Description

octav is the API server that does the basic CRUD operations on data.
No authentication is done on this component.

All basic text data is expected to be registered using English, but
each data component may have elements that can be localized.
You can register such data by providing data with keys such as `name#ja`

Localized data can be retrieved by providing the `lang` key to each of
the `Lookup*` endpoints. For example, to fetch sessions with localized
title and abstract, you can issue a request like the following:

```
http://******/v1/session/lookup?id=*******&lang=ja
```

As you can see, endpoints that fetch data are normally represented using
GET endpoints, with data encoded in query strings.

For other endpoints that require register, update, or delete data are
represented using POST endpoints. In this case the data should be
encoded as JSON text.

# Code generation directions

## If you edited spec/v1/api.json

```
make buildspec
```

This will regenerate 

* octav/client/client.go
* octav/validator/validator.go
* octav/octav_hsup.go

If you add more endpoints, you need to write additional `doXXXX` handlers in handlers.go

## If you edited interface.go or db/interface.go

Technically, you do NOT need to do this every time, but if you want to be safe, whenever you touch these files you should regenerate auto-generated files

To regenerate files, run:

```
make generate
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

