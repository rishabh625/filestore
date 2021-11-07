# Filestore

Distributed File Store Application Consist of API Server to handle file operations and command line tool to do operations (store named binary is present inside docker container) , Hash of the file content is generated at client side and is stored in Redis at server side

To use command line set below env variable
```
API_HOST=localhost //host where API server is running
API_PORT=5000 //port on which API server is running

```
### Steps to Build cli client
```
go build -tags static_all -a -installsuffix cgo -ldflags '-extldflags "-static"' -o 'store' client/main.go

```
and Use ./store help 

To use server set below env variable and run server/main.go
```
REDIS_HOST=localhost //host where redis server is running
REDIS_PORT=5000 //port on which redis server is running
```

### Steps to Use on local machine via docker-compose

1) git clone this repo
2) Create folder as redis-data
3) docker-compose up --build -d
4) docker exec -it filestore_app_1 bash
5) ./store --help // to access all commands
