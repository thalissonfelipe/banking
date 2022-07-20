# Banking

Banking is a go project using [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). This project is a technical challenge proposed by the Stone company.

Banking is a REST API that allows users to create accounts and to perform transactions.

## Technologies and libraries used in this project

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [chi](https://github.com/go-chi/chi)
- [jwt](https://github.com/golang-jwt/jwt)
- [pgxpool](https://github.com/jackc/pgx)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [brdoc](https://github.com/Nhanderu/brdoc)
- [envconfig](https://github.com/kelseyhightower/envconfig)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20210513164829-c07d793c2f9a/bcrypt)
- [cli](https://github.com/urfave/cli)
- [dockertest](https://github.com/ory/dockertest)
- [swag](https://github.com/swaggo/swag)
- [gRPC](https://grpc.io/)
- [otel](https://opentelemetry.io/)
- [zap](https://github.com/uber-go/zap)
- [golangci](https://github.com/golangci/golangci-lint)
- [buf.build](https://buf.build/)

## Getting Started

To get started quickly, you must have `docker`, `docker-compose` and `make` installed in your machine.

1. If you want to configure your own API and database settings, create an .env file in the root folder and fill it out following the .env.example file as a template. Be aware that if you change the default portgres configuration, you may need to configure something manually inside the Docker container.

2. Run the following command if you want to start the application inside docker container:
    ```shell
    make dev-docker
    ```

3. (Optional) Run the following command if you want to start the application locally. This option requires Go. Also, since your application will run locally and the postgres database inside docker, you need to set the database variables directly in the docker-compose.yml file, or export the variables using the keyword export, or fill the .env file as stated on step one (recommended).
    ```shell   
    make dev-local
    ```

The server will be listening on `localhost:5000`.

## Documentation

Swagger URL: http://localhost:5000/swagger/index.html

## Tests

```shell
make test
```
