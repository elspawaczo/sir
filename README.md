# SOON\_ Instance Registry (`sir`)

A simple Go web application which registers servers and assigning them a new hostname. This is intended
for use within a Private network and therefor has no authentication. Redis is used as a store for a pool
of available hostnames and taken hostnames.

## Importing

To import data into Redis you must have a line delimited text file containing names you
wish to add into Redis, for example:

```
foo
bar
baz
```

You can then run the following import command:

```
sir import --path /path/to/file.txt
```

This will import the file into a Redis set running on localhost.

## Supported Endpoints

### `GET /`

This will return information about the current registry.

``` http
GET / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8000
User-Agent: HTTPie/0.9.2
```

``` http
HTTP/1.1 200 OK
Content-Length: 49
Content-Type: application/json
Date: Mon, 06 Jul 2015 16:36:17 GMT

{
    "available": 1045,
    "remaining": "1.88%",
    "taken": 20
}
```

### `POST /`

Allocates a hostname from the pool and saves instance id to that allocated hostname.

``` http
POST / HTTP/1.1
Accept: application/json
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 14
Content-Type: application/json
Host: localhost:8000
User-Agent: HTTPie/0.9.2

{
    "instance_id": "i-abc1234"
    "private_ip": "10.0.0.1"
}
```

``` http
HTTP/1.1 200 OK
Content-Length: 23
Content-Type: application/json
Date: Mon, 06 Jul 2015 16:37:01 GMT

{
    "name": "triggerhappy"
}
```

### `GET /:instance_id`

Will return the hostname allocated to the instance and other related information.

``` http
GET /i-abc1234 HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8000
User-Agent: HTTPie/0.9.2
```

``` http
HTTP/1.1 200 OK
Content-Length: 49
Content-Type: application/json
Date: Mon, 06 Jul 2015 16:36:17 GMT

{
    "instance_id": "i-abc1234",
    "name": triggerhappy,
    "private_ip": "10.0.0.1"
}
```

### `DELETE /:instance_id`

Will remove the hostname from the taken pool allowing it to be reallocate.

``` http
DELETE /i-abc1234 HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 0
Host: localhost:8000
User-Agent: HTTPie/0.9.2
```

``` http
HTTP/1.1 204 No Content
Content-Type: application/json
Date: Mon, 06 Jul 2015 16:37:42 GMT
```

## Developing

To add new functionality to `sir` you will need:

* Go 1.4+
* Go Package Manager (`gpm`)

### WorkSpace

Set up your workspace by creating a the directory layout `Go` expects. All examples
will assume you are working from `~/Development`.

```
mkdir -p ~/Development/sir/src/github.com/thisissoon/sir
```

Now cd into that diectory and clone the repository:

```
cd ~/Development/sir/src/github.com/thisissoon/sir
git clone git@github.com:thisissoon/sir.git .
```

Now set your go path to be the top level directory of your workspace and update your
`bin` path to include your workspace.

```
export $GOPATH=~/Development/sir
export PATH=$PATH:$GOPATH/bin
```

You can now install the dependencies:

```
gpm install
```

Once everything is installed you can now build / install the project:

```
go build ./...
```

A `sir` command will now be on your `$PATH`

```
sir --help
which sir
~/Development/sir/bin/sir
```

### Third Parties

`sir` uses a couple of third parties:

* Cobra - A POSIX command line application builder: github.com/spf13/cobra
* Rredis - A Go Redis Client: gopkg.in/redis.v3
* Goji - A Go Micro Web Framework: github.com/zenazn/goji
* Validator - A simple Go Validator: gopkg.in/validator.v2

These will be installed by `gpm`
