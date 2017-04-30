NAME            := pulpe
PACKAGES        := $(shell glide novendor)

.PHONY: all build $(NAME) deps install gen test testrace

all: $(NAME)

$(NAME):
	go install ./cmd/$@

buildstatic:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/pulpe

build: clean buildstatic dist docker

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

dist:
	cd web/ && yarn run build:prod

docker:
	docker build -t blankrobot/pulpe .

clean:
	rm -fr dist/
