# omnilint-server

Send your code to this server to have it checked for errors, warnings, etc.

This is my first half-serious use of Google's Go, so please tell me where the
code can improve, etc.

## Language Support

### PHP

- set `Content-Type` HTTP header to `application/x-php`

Code is scanned with:

- the basic parser: `php -l`
- [PHP_CodeSniffer](https://github.com/squizlabs/PHP_CodeSniffer): `phpcs`

#### Options (Query String Arguments)

- `phpcs.standard`: defaults to "PSR2"

## HTTP routes

### POST *

- any path, path doesn't matter
- expects `Content-Type` HTTP header to be set
- expects the code to be scanned to occupy the entire HTTP request body

## Development

- this project uses [godep](https://github.com/tools/godep) to manage
  package dependencies

## Building the Docker image

- install [gox](https://github.com/mitchellh/gox) and set it up for 64-bit Linux

```shell
gox -osarch="linux/amd64"
cp $GOPATH/bin/omnilint-server-xc/snapshot/omnilint-server_linux_amd64.tar.gz .
docker build -t IMAGE .
rm omnilint-server_linux_amd64.tar.gz
```

## Running the Docker image

```shell
docker run -t -i -p 3000:3000 IMAGE
```

- use `-e` to set environment variables as required

## Environment Variables

- `PORT` and `HOST` to control interface binding (default is `localhost:3000`)

- `NEWRELIC_LICENSE` and `NEWRELIC_NAME` for integration with [New Relic](http://newrelic.com/)

- `CORS_ORIGINS` (with comma-separated origins) to enable [CORS HTTP headers](http://www.html5rocks.com/en/tutorials/cors/)
