# omnilint-server

Send your code to this server to have it checked for errors, warnings, etc.

## Language Support

### PHP

- set `Content-Type` HTTP header to `application/x-php`

Code is scanned with:

- the basic parser: `php -l`

## HTTP routes

### POST *

- any path, path doesn't matter
- expects `Content-Type` HTTP header to be set
- expects the code to be scanned to occupy the entire HTTP request body
