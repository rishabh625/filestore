# Filestore

Distributed File Store Application Consist of API Server to handle file operations and command line tool to do operations

To use command line set below env variable
```
API_HOST=localhost //host where API server is running
API_PORT=5000 //port on which API server is running

```

### Steps to Use on local machine

1) git clone this repo
2) docker-compose up --build -d
3) docker exec -it filestore_app_1 bash
4) ./store --help // to access all commands
