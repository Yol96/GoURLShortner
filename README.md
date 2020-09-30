# GoUrlShortner

A URL Link shortner API written in Go.
The purpose of the API is to shorten URL links via REST HTTP request.


### Install and run the server
Using Makefile:
```sh
$ make
$ ./apiserver
```
Using Docker-compose:
```sh
$ docker-compose build
$ docker-compose up
```

Also u can add .env file in root folder, for example:
```sh
APP_SERVER_PORT=":8080"
APP_LOGRUS_LEVEL="debug"
APP_REDIS_ADDRESS="redis:6379"
APP_REDIS_PASSWORD=""
APP_REDIS_DB=0
```


### POST
The POST request should include the address and expiration_time parameters as follows:
```sh
{
    "address":"https://www.google.com/",
    "expiration_time":0
}
```

![](/screens/post_new.jpg)

### GET
You can get information about URL by a GET request to http://localhost:8080/info?link=shortlink
![](/screens/get_info.jpg)

### REDIRECT
You can consume a short URL by issuing a GET request to http://localhost:8080/shortlink
![](/screens/redirect.jpg)

### TODO:
- [ ] Add validations
- [x] Add app logging
- [x] Add errors handling(error.go)
- [ ] Add possibility to create short links with given name