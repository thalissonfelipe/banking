# Banking

Banking is a go project using [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). This project is a challenge proposed by the Stone company.

Banking is a REST API that allows users to create accounts and to perform transactions.

## Technologies

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)

## Libraries

- [gorilla/mux](https://github.com/gorilla/mux)
- [jwt-go](https://github.com/dgrijalva/jwt-go)
- [uuid](https://github.com/google/uuid)
- [pgx](https://github.com/jackc/pgx)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [brdoc](https://github.com/Nhanderu/brdoc)
- [envconfig](https://github.com/kelseyhightower/envconfig)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20210513164829-c07d793c2f9a/bcrypt)
- [cli](https://github.com/urfave/cli)
- [testify](https://github.com/stretchr/testify)

## Getting Started

To get started quickly, you must have `Docker` and `make` installed in your machine to start the application easily.

1. If you want to configure your own API and database settings, create an .env file in the root folder and fill it out following the .env.example file as a template. Be aware that if you change the default portgres configuration, you may need to configure something manually within the Docker container.

2. Run the following command if you want to start the application inside docker container:
    ```shell
    make dev-docker
    ```

3. (Optional) Run the following command if you want to start the application locally. This options requires Go.
    ```shell
    make dev-local
    ```
The server will be listening on `0.0.0.0:5000`.

## Endpoints

- **GET /accounts**
    - List all accounts.
- **GET /accounts/{id}/balance.**
    - Gets the account balance by ID, if exists.
- **POST /accounts**
    - Creates a new account.
- **POST /login**
    - If login succeeds, a jwt token will be returned and must be used on endpoints that involve transfers.
- **GET /transfers**
    - Lists all transfers. The user must be authenticated to use this endpoint.
- **POST /transfers**
    - Performs a new transfer between accounts. The user must be authenticated to use this endpoint.

## Tests

To run the tests:

```shell
make test
```

## TODO

- [ ] Improve endpoints documentation.
- [ ] Add swagger.
- [ ] Add logger.
- [ ] Add integration tests.
- [ ] Add coverage tests.
- [ ] Add unit tests for repository methods.
- [ ] Fix json responses (the names are capitalized).
- [ ] Add endpoint to deposit money (the default balance is 0, so it's not possible to perform transaction since the user does not have funds).
- [ ] Add Github Actions.
- [ ] Add prefix endpoint (must start with /api/v1/ or something like that).
- [ ] Add common pkg to centralize some functions, errors, constants...
- [ ] Change id saved in the bank to be a UUID type (currently the id is being saved as a string).
- [ ] Move creation datetime do PostgreSQL.
- [ ] Add CORS middleware.
