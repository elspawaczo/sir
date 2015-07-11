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
    "hostname": triggerhappy,
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
