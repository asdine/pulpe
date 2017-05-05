# Pulpe

[![GoDoc](https://godoc.org/github.com/blankrobot/pulpe?status.svg)](https://godoc.org/github.com/blankrobot/pulpe)
[![Build Status](https://travis-ci.org/blankrobot/pulpe.svg)](https://travis-ci.org/blankrobot/pulpe)

Pulpe is an open source web application for managing content.

*Note: Work in progress*

## Build from source

### Requirements

- [Go](https://golang.org/)
- [Node.js 7.x](https://nodejs.org)
- [Glide](https://github.com/Masterminds/glide)
- [Yarn](https://yarnpkg.com/)
- [Docker](https://www.docker.com/)

### Install dependencies

Go dependencies:

```sh
make install
```

Front dependencies:

```sh
cd web && yarn install
```

## Dev

Compile the Go server

```sh
make
```

Run the database

```sh
docker-compose up -d
```

Run the dev server

```sh
cd web && yarn run dev
```
