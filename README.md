# user-api-golang

A CRUD API template for perform operations on users. Written in golang.

## Components

- Server: Golang with [go-chi](https://github.com/go-chi/chi) as mux and [zap](https://github.com/uber-go/zap) for logging
- JWT auth with [go-chi/jwtauth](https://github.com/go-chi/jwtauth)
- Posgres with [go-pg](https://github.com/go-pg/pg) library

## Running instructions

1. Make sure docker is installed and then run `docker compose up --build` from base of the repo.
2. Server should be up on `localhost:8000/`

## Scenarios

### Signup

```
curl -X POST http://localhost:8000/signup -v\
   -H 'Content-Type: application/json' \
   -d '{"email":"test@example.com","password":"somepass","firstName":"John","lastname":"Doe"}'
```

### Login

```
curl -X POST http://localhost:8000/login -v\
   -H 'Content-Type: application/json' \
   -d '{"email":"test@example.com","password":"somepass"}'
```

### Update user information

```
curl -X PUT http://localhost:8000/users -v\
   -H 'Authorization: BEARER <token>' \
    -H 'Content-Type: application/json' \
    -d '{"firstName":"Johny","lastname":"Doesnt"}'
```
