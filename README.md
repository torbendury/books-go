# 📚 books-go

[![Go Report Card](https://goreportcard.com/badge/github.com/torbendury/books-go)](https://goreportcard.com/report/github.com/torbendury/books-go)

A small playground application which consists of a REST API (crud actions) that allows you to create, read, update and delete books.

- [📚 books-go](#-books-go)
  - [🧬 Repo structure](#-repo-structure)
  - [🖼️ Overview](#️-overview)
  - [👷 Usage](#-usage)
    - [🏃 Running](#-running)
  - [✔️ TODOs](#️-todos)
  - [📷 Generating puml](#-generating-puml)

## 🧬 Repo structure

```txt
.
├── README.md
├── api               # contains the actual server logic. CRUD routes, middleware registration, ...
├── data              # models for our backend data
├── hack              # HTTP requests and docker-compose file for spinning up a local DB
├── cmd/main.go       # firestarter for application. holds some config and does nothing else than starting up.
├── storage           # storage interface to keep interchangeable between in-memory and other storages
└── utilities         # unused, might come in handy later.
```

## 🖼️ Overview

![PlantUML graphic](https://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/torbendury/books-go/main/docs/graph.puml)

## 👷 Usage

### 🏃 Running

Run `make run`. The server will start up and be reachable at [`http://localhost:3000`](http://localhost:3000).

See [`test.http`](hack/test.http) for available API endpoints.

If you want to use PostgreSQL as a backend and don't have any at hand, you can use the [docker-compose](hack/docker-compose.yml) file to spin one up, or use the Makefile directly:

```bash
make startpsql
make runpg

# when you're finished, stop the PSQL and clean up
make stoppsql
```

**NOTE:** The DB init scripts only run on the first time. To run them again and effectively clean up your database, you will want to prune the docker volumes after a shutdown. See also [notes](hack/NOTES.md).

## ✔️ TODOs

See [TODO](TODO).

## 📷 Generating puml

As seen above, you can generate PlantUML code from the repo. Run `make puml` for this.

Note that until [this issue](https://github.com/bykof/go-plantuml/issues/27) is resolved, you have to manually adjust the PlantUML layout with `left to right direction`.
