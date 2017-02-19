# Pulpe

[![GoDoc](https://godoc.org/github.com/blankrobot/pulpe?status.svg)](https://godoc.org/github.com/blankrobot/pulpe)
[![Build Status](https://travis-ci.org/blankrobot/pulpe.svg)](https://travis-ci.org/blankrobot/pulpe)

Pulpe is an open source web application for managing content.

*Note: Work in progress*

## Build from source

### Requirements

- [Go](https://golang.org/)
- [Node.js 6.x](https://nodejs.org)
- [Glide](https://github.com/Masterminds/glide)
- [Yarn](https://yarnpkg.com/)
- [Docker](https://www.docker.com/)

### Install dependencies

Go dependencies:

```sh
$ make install
```

Front dependencies:

```sh
$ cd web && yarn install
```

## Dev

Compile the Go server
```sh
$ make
```

Run the database
```sh
$ docker-compose up -d
```

Run the dev server
```sh
$ cd web && yarn run dev
```

// Board
pulpe.io/asdine/recettes-de-cuisine

// List
pulpe.io/asdine/recettes-de-cuisine/desserts

// Card
pulpe.io/asdine/recettes-de-cuisine/desserts/tarte-aux-fraises

// Direct access
pulpe.io/b/58a822d57aa61b3ff92ed2d8
pulpe.io/l/58a822d97aa61b3ff92ed2d9
pulpe.io/c/58a822df7aa61b3ff92ed2da

----------------------------------------------

// API boards
POST    pulpe.io/v1/boards
GET     pulpe.io/v1/boards/asdine/recettes-de-cuisine
        pulpe.io/v1/boards/58a822d57aa61b3ff92ed2d8
PATCH   pulpe.io/v1/boards/58a822d57aa61b3ff92ed2d8
DELETE  pulpe.io/v1/boards/58a822d57aa61b3ff92ed2d8

// API lists
POST    pulpe.io/v1/boards/asdine/recettes-de-cuisine/lists
GET     pulpe.io/v1/lists/58a822d97aa61b3ff92ed2d9
PATCH   pulpe.io/v1/lists/58a822d97aa61b3ff92ed2d9
DELETE  pulpe.io/v1/lists/58a822d97aa61b3ff92ed2d9

// API cards
POST    pulpe.io/v1/lists/58a822d97aa61b3ff92ed2d9/cards
GET     pulpe.io/v1/cards/58a822df7aa61b3ff92ed2da
PATCH   pulpe.io/v1/cards/58a822df7aa61b3ff92ed2da
DELETE  pulpe.io/v1/cards/58a822df7aa61b3ff92ed2da

----------------------------------------------
