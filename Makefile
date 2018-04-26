TSERVER:=marge-backend
DOCKER_REPO:=albertocsm

.PHONY: all clean proto client server build docker docker-publish docker-tag

all: clean proto build docker docker-tag docker-publish

proto:
	protoc --go_out=plugins=grpc:${GOPATH}/src grpc/*.proto

clean:
	rm -rf server/${TSERVER}

build:
ifndef TGT_LINUX
	cd server && CGO_ENABLED=0 go build -o ${TSERVER} *.go
else
	cd server && CGO_ENABLED=0 GOOS=linux go build -o ${TSERVER} *.go
endif


docker:
	cd server && docker build --tag=${TSERVER} .

docker-publish:
	docker tag ${TSERVER} albertocsm/${TSERVER}
	docker push albertocsm/${TSERVER}

docker-push: docker docker-publish

