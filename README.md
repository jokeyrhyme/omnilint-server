# omnilint-server

Send your code to this server to have it checked for errors, warnings, etc.

## Language Support

### PHP

- set `Content-Type` HTTP header to `application/x-php`

Code is scanned with:

- the basic parser: `php -l`
- [PHP_CodeSniffer](https://github.com/squizlabs/PHP_CodeSniffer): `phpcs`

## HTTP routes

### POST *

- any path, path doesn't matter
- expects `Content-Type` HTTP header to be set
- expects the code to be scanned to occupy the entire HTTP request body

## Building the Docker image

- install [goxc](https://github.com/laher/goxc) and set it up for 64-bit Linux

```shell
goxc -bc="linux"
cp $GOPATH/bin/omnilint-server-xc/snapshot/omnilint-server_linux_amd64.tar.gz .
docker build -t IMAGE .
rm omnilint-server_linux_amd64.tar.gz
```

## Running the Docker image

```shell
docker run -t -i -p 3000:3000 IMAGE
```

- use `-e` to set NEWRELIC_LICENSE and NEWRELIC_NAME environment variables (if
  you want to integrate with New Relic)
