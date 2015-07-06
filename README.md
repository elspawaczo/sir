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

Will allocate a hostname and return it back.

``` http
POST / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 0
Host: localhost:8000
User-Agent: HTTPie/0.9.2
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

### `DELETE /:hostname`

Will remove the hostname from the taken pool allowing it to be reallocate.

``` http
DELETE /triggerhappy HTTP/1.1
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
