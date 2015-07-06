# SOON\_ Instance Registry (`sir`)

A simple Go web application which registers servers and assigning them a new hostname. This is intended
for use within a Private network and therefor has no authentication. Redis is used as a store for a pool
of available hostnames and taken hostnames.

## Supported Endpoints

### `GET /`

This will return information about the current registry.

### `POST /`

Will allocate a hostname and return it back.

### `GET /:hostname`

Will return information about an allocated hostname.

### `DELETE /:hostname`

Will remove the hostname from the taken pool allowing it to be reallocate.
