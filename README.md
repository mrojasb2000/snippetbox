

## Customizing HTTP Headers

### Check whether the request is using POST or not.

```sh
$ curl -i  http://localhost:3000/snippet/create
HTTP/1.1 405 Method Not Allowed
Date: Thu, 11 Aug 2022 21:28:33 GMT
Content-Length: 18
Content-Type: text/plain; charset=utf-8
```

```sh
curl -i -X POST  http://localhost:3000/snippet/create
HTTP/1.1 200 OK
Date: Thu, 11 Aug 2022 21:28:53 GMT
Content-Length: 23
Content-Type: text/plain; charset=utf-8

Create a new snippet...
```


### Improvement we can make is to include an Allow: POST header

```sh
curl -i  http://localhost:3000/snippet/create
HTTP/1.1 405 Method Not Allowed
Allow: POST
Date: Thu, 11 Aug 2022 21:43:37 GMT
Content-Length: 18
Content-Type: text/plain; charset=utf-8
```