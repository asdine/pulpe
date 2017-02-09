NAME            := pulpe
PACKAGES        := $(shell glide novendor)
TEST_MONGO_URI  := "mongodb://localhost:27017"

.PHONY: all build $(NAME) deps install gen test testrace

all: build

build: $(NAME)

$(NAME):
	go install ./cmd/$@

deps:
	glide up

install:
	glide install

gen:
	go generate $(PACKAGES)

test:
	MONGO_URI=$(TEST_MONGO_URI) go test -v -cover $(PACKAGES)

testrace:
	MONGO_URI=$(TEST_MONGO_URI) go test -v -race -cover $(PACKAGES)
