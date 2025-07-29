# GRAFIT PRODUCE
## Service
### Transaction Notification
Technical test base on studi case

## Features
- Credit Plus
- Check Status

## Endpoint
list endpoint you can see as json file in [endpoint auth service](cmd/api/docs/swagger.json), end then copy that json and you can see in [editor swagger](https://editor.swagger.io/)

### Generate endpoint documentation
while you want to self generate endpoint documentation you can following this command:
```bash
$ swag init --generalInfo ./cmd/api/main.go --output ./cmd/api/docs
```

## How to run
### Specification
- database: postgres
- language: golang:1.24.5
- environtment: file [.env](../configs/file.env) copy and paste in auth project as `.env`
- migrate: make sure you migration is already setup on global setup
- other: golang-migrate, swaggo
- storage: minio

#### auto install dependecies golang, and golang-migrate just run ```./install-dep.sh``` or setup dependencies using docker just run ```docker-compose up -Vd dbpgsql adminpgsql```

### How to run
1.  clone this project
```bash
$ git clone https://github.com/Grafiters/technicalTest.git
```

2.  install all dephendencies go.mod
```bash
$ go mod tidy
```

3.  running migration
```bash
$ sudo chmod +x script/task-list.sh
$ ./script/task-list.sh db:create
$ ./script/task-list.sh db:migrate
$ ./script/task-list.sh db:seed
```

4.  running project
    - running http endoint you can use:
    ```bash
    $ go run cmd/api/main.go
    ```

5.  running worker
    - running http endoint you can use:
    ```bash
    $ go run cmd/workers/main.go
    ```

## How To Deploy
### Deploy using docker
Deployment launchpad service can using dockerize with build image and run with docker-compose, follow step deployment in below:
1.  make sure the project has been cloning to your local machine computer
2.  build the project using docker, exec the comand line below:
```bash
$ docker build -t credit-plu:0.0.1 .
```

3.  migrate all migratin table, exec the command line below:
```bash
$ docker-compose run --rm api sh -C "script/task-list.sh db:create && script/task-list.sh db:migrate && script/task-list.sh db:seed"
```
4.  then run all service like api, and worker, you can see the file docker-compose on [docker-compose.yml](./docker-compose.yml)
```bash
$ docker-compose up -Vd
```

# Contact
For question, please contact telegram [alone](https://t.me/Grafiters)
