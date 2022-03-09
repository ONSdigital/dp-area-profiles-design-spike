POSTGRES_USER=
POSTGRES_PASSWORD=
DATABASE_NAME=

.PHONY: build
build:
	@echo "building"
	go build -o api

.PHONY: fresh
fresh: build
	@echo "starting api with a clean database"
	./api -u=${POSTGRES_USER} -p=${POSTGRES_PASSWORD} -db=${DATABASE_NAME} -drop=true

.PHONY: run
run: build
	@echo "starting api with database as is"
	./api -u=${POSTGRES_USER} -p=${POSTGRES_PASSWORD} -db=${DATABASE_NAME} -drop=false

.PHONY: compose
compose:
	POSTGRES_USER=${POSTRES_USER} POSTGRES_PASSWORD=${POSTGRES_PASSWORD} docker-compose up