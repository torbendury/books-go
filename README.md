# 📚 books-go

[![Go Report Card](https://goreportcard.com/badge/github.com/torbendury/books-go)](https://goreportcard.com/report/github.com/torbendury/books-go)

A small playground application which consists of a REST API (crud actions) that allows you to create, read, update and delete books.

- [📚 books-go](#-books-go)
  - [🧬 Repo structure](#-repo-structure)
  - [👷 Usage](#-usage)
    - [🏃 Running](#-running)
  - [✔️ TODOs](#️-todos)
  - [🖼️ Overview](#️-overview)

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

## 👷 Usage

### 🏃 Running

Run `make compose`. The server will start up and be reachable at [`http://localhost:3000`](http://localhost:3000).

See [`test.http`](hack/test.http) for available API endpoints which are ready for usage when you spin up the application.

## ✔️ TODOs

See [TODO](TODO).

## 🖼️ Overview

![PlantUML graphic](https://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/torbendury/books-go/main/docs/graph.puml)
