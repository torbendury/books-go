# books-go

A small playground application which consists of a REST API (crud actions) that allows you to create, read, update and delete books.

- [books-go](#books-go)
  - [Repo structure](#repo-structure)
  - [Overview](#overview)
  - [Usage](#usage)
    - [Running](#running)
  - [TODOs](#todos)
  - [Generating puml](#generating-puml)

## Repo structure

```txt
.
├── README.md
├── api               # contains the actual server logic. CRUD routes, middleware registration, ...
├── data              # models for our backend data
├── go.mod
├── go.sum
├── cmd/main.go       # firestarter for application. holds some config and does nothing else than starting up.
├── storage           # storage interface to keep interchangeable between in-memory and other storages
├── test.http
└── utilities         # unused yet, might come in handy later.
```

## Overview

![tech overview](https://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/torbendury/books-go/main/docs/graph.puml)

## Usage

### Running

After installing dependencies using `go get`, you should be able to run the project with `go run cmd/main.go`. The server will start up and be reachable at [`http://localhost:3000`](http://localhost:3000).

See [`test.http`](test.http) for available API endpoints.

## TODOs

- Implement PostgresqlStorage
- Fix updating books
- Write docs!
- Go and see if all those dependencies are "needed" for this small project or can be cleaned up

## Generating puml

Run: `go-plantuml generate -d . -r`
