# sequence-technical-test

## Requirements

To run the images:

- Docker version 28 or greater with Docker Compose

Or to run the application locally:

- Golang 1.25+
- Postgres 17+

### Environment Variables

All environment variables are available in the `.env.example` file, you can copy them to a `.env` file to test the app locally

## Running

The recommended way to run this application is through Docker Compose.

### Using Docker Compose

```shell
docker compose up --build -d 
```

This will create a postgres container, a pgAdmin app to interact direclty with the database via web interface, a container that automatically runs all the migrations and the Sequence API

### Using Docker

To run it with just docker, you'll need to create a postgres container first:

```shell
docker run --name db -d -p 5432:5432 \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=postgres \
    postgres:17.6-alpine3.22
```

Then, create a `.env` file in the root of the project with the postgres credentials, this is required to run the migrations in the database.

After that, run the migrations:

```shell
make folder=postgres migrate
make folder=sequencemailbox migrate
```

With the migrations in place, the database is ready to serve the application.

Build and run the app's container:

```shell
docker build -t sequence-technical-test .

docker run --name sequence-technical-test -d -p 8000:8000 \
    -e APP_PORT=8000 \
    -e DB_HOST=db.local \
    -e DB_PORT=5432 \
    -e DB_USER=postgres \
    -e DB_PASSWORD=postgres \
    -e DB_NAME=postgres \
    -e DB_MAX_CONNECTIONS=10 \
    -e DB_MIN_CONNECTIONS=1 \
    -e DB_MAX_CONN_IDLE_TIME=30 \
    sequence-technical-test
```

### Using Golang Locally

This application depends on Postgres, Golang and some environment variables described at the Requirements section, so you'll need to install/set them first.

Then, build and run the sequence-technical-test app:

```shell
make clear

make build

make run-bin
```

After that, the app will be available at the specified port or, by default, port 8000.

### [Optional] Install MockGen

Mockgen is a golang tool used to create mocks for unit testing based on application's interfaces.

It is not required to just run the application, but it is an important tool to create new unit tests.

First, make sure that your golang bin folder is in yout PATH. Then, run the install_tools command:

```shell
make install_tools
```

This will install mockgen and it should be ready to use. Check if its installed succefully with:

```shell
mockgen -version
```
