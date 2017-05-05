NAME            := pulpe
PACKAGES        := $(shell glide novendor)

.PHONY: all build $(NAME) deps install gen test testrace

all: $(NAME)

$(NAME):
	go install ./cmd/$@

build: clean $(NAME) dist

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

docker: clean
	mkdir -p dist/
	docker build -t blankrobot/pulpe-web-builder -f ./web/Dockerfile ./web
	docker run -v $(PWD)/dist:/dist blankrobot/pulpe-web-builder
	docker build -t blankrobot/pulpe-builder -f Dockerfile.build .
	docker run -v $(PWD)/dist:/dist blankrobot/pulpe-builder
	docker build -t blankrobot/pulpe .

clean:
	rm -fr dist/
	rm -f $(NAME)
