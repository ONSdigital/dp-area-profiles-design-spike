.PHONY: build
build:
	go build -o poc

## Start the API and drop any existing data/tables and recreate the schema.
.PHONY: debug
debug: build
	./poc import -f="load/ex1.csv" -d=true

.PHONY: query
query: build
	./poc query

## Start the API and drop any existing data/tables and recreate the schema.
.PHONY: import
import: build
	./poc -u=${POSTGRES_USER} -p=${POSTGRES_PASSWORD} -db=${DATABASE_NAME} -cmd=import

## Start the API and drop any existing data/tables and recreate the schema.
.PHONY: clean
clean: build
	./poc -u=${POSTGRES_USER} -p=${POSTGRES_PASSWORD} -db=${DATABASE_NAME} -cmd=clean

## Start the API retaining the current database state.
.PHONY: run
run: build
	./poc -u=${POSTGRES_USER} -p=${POSTGRES_PASSWORD} -db=${DATABASE_NAME} -drop=false

## Runs a Postgres Docker container for the app to connect to.
.PHONY: compose
compose:
	POSTGRES_USER=${POSTGRES_USER} POSTGRES_PASSWORD=${POSTGRES_PASSWORD} docker-compose up