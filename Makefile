NAME            := pulpe
PACKAGES        := $(shell glide novendor)

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
	 go test -v -cover $(PACKAGES)

testrace:
	go test -v -race -cover $(PACKAGES)
