POSTGRES_USER=postgres ## container default (safe to commit)
POSTGRES_PASSWORD=mysecretpassword ## container default (safe to commit)
DATABASE_NAME=area_profiles

.PHONY: build
build:
	go build -o api

## Start the API and drop any existing data/tables and recreate the schema.
.PHONY: fresh
fresh: build
	./api -u=${POSTGRES_USER} -p=${POSTGRES_PASSWORD} -db=${DATABASE_NAME} -drop=true

## Start the API retaining the current database state.
.PHONY: run
run: build
	./api -u=${POSTGRES_USER} -p=${POSTGRES_PASSWORD} -db=${DATABASE_NAME} -drop=false

## Runs a Postgres Docker container for the app to connect to.
.PHONY: compose
compose:
	POSTGRES_USER=${POSTRES_USER} POSTGRES_PASSWORD=${POSTGRES_PASSWORD} docker-compose up