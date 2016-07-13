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

# Data structure

## Conference Series

"builderscon" is a series. "builderscon tokyo 2016" is a conference instance in that series.

Only system administrators can create conference series.

## Conferences

Conferences have title, sub title, series_id and slug. The slug is used to create
URL (http://conf.builderscon.io/*slug*)

Conferences can have multiple dates. Most conferences will be 1 day only,
but sometimes they span multiple days. Use `AddConferenceDates` to
add dates to a conference.

Conferences must have at least 1 administrator associated, but can have
multiple administrators. The account that created the conference is
automatically registered as administrator. To add more administrators,
use `AddAdminisrator`

Conferences may only be created by conference series administrators, but may be edited by conference administrators.

## ConferenceDates

Conference dates take the following format:

```
YYYY-MM-DD[HH:MM-HH:MM]
YYYY-MM-DD[HH:MM]
YYYY-MM-DD
```

When you specify the time, it's in one of these forms `[OpenTime-CloseTime]`, `[OpenTime]`

You can add a new date using `AddConferenceDates` API call, and remove them using `DeleteConferenceDates`.

## Conference Administrators

You may register as many conference administrators as you like. The original creator
of the conference resource is automatically associated as the administrator.

You can add more administrators by using the `AddConferenceAdmin` API call, and remove
them using `DeleteConferenceAdmin`
 
## Venues

Venues can contain multiple rooms.

Venues can be created by anybody who is an administrator of any conference.

## Rooms

Rooms belong to a venue.

Rooms can be created by anybody who is an administrator of any conference.

## Sessions

Sessions can be regeistered to be presented at a room.

Sessions for a particular conference can only be created by administrators of that conference.

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

## If you edited model/interface.go or db/interface.go

Technically, you do NOT need to do this every time, but if you want to be safe, whenever you touch these files you should regenerate auto-generated files

To regenerate files, run:

```
make generate
```

If you add transport types (e.g. `CreateSessionRequest`, `UpdateSessionRequest`, etc) make sure to include a flag in the comment, or otherwise the code generation tools do not pick them up:

```go
// +transport
type CreateFooRequest struct {
   ...
}
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
make initdb
```

## 1. Create a database

```
make initdb
```

## 2. Run

```
cd octav
OCTAV_TEST_DSN='root:@:/octav?parseTime=true' go test .
```

Currently, the tests are still failing.

# Deploying to GKE

There are LOTS of things that needs to be automated. Please HELP!

## Get CloudSQL Credentials

Ask an administrator for them.

Once you have them, you will be able to access the CloudSQL from your
terminal by typing:

```
make cloudsql
```

## Build and Deploy the container

```
cd gke/containers/apiserver
make docker
make publish # Remember the tag that is printed
make deploy TAG=XXXXX
```

`make docker` builds octav, and creates a docker container.

`make publish` will display a line like the following:

```
Publishing [ asia.gcr.io/builderscon-1248/apiserver:20160317.212833 ]
```

You need to know the tag (20160317.212833) to deploy. You can also find it
in the GCP console page.

`make deploy` deploys the specified tag using `kubectl rolling-update`

## kubectl

Random `kubectl` commands that are good to know:

| Command | Description |
|:--------|:------------|
| kubectl get pods | List all pods |
| kubectl get rc   | List all replication controllers |
| kubectl get service | List all services |
| kubectl delete pod -l name=apiserver | Delete all `apiserver` pods. If there's a replication controller, new instances will be brought up |
| kubectl delete rc -l name=apiserver | Delete all `apiserver` replication controllers. This also kills associated pods |
| kubectl logs \[-f\] \[pod name\] | Show standard output logs. -f does the equivalent of `tail` |
| kubectl describe \[name\] | name can be any name like pod name, rc name, service name, etc. Gives detail information about the resource |
| kubectl exec -it \[pod name\] /bin/sh | Open a shell session to the named pod. Note: if you have multiple containers in a pod, you need add the container name |
