# dp-area-profiles-design-spike

POC API demonstrating a proposed relational database schema to model Area Profiles & key statistics. The POC loads some
test data for area profiles & key stats which can then be queried via some API endpoints.

:warning: This is a POC, not production ready code and its for illustrative purposes only. :warning:

### Getting started

Get the code:

  ```bash
  git clone https://github.com/ONSdigital/dp-area-profiles-design-spike
  ```

The POC is backed by a Postgres DB which needs to be set up before running the app. Ensure Docker is running then from
the project root dir run:

  ```
  make compose
  ```

Open another terminal and run the following to connect to Postgres:

```bash
docker exec -it dp-area-profiles-design-spike_postgres_1 psql -U postgres
```

** :warning: The container name my vary. Use `docker ps` to get the name of your conatiner and replace
`dp-area-profiles-design-spike_postgres_1` as required.

Create a database for the POC to connect to:

  ```
  CREATE DATABASE area_profiles;
  ```

### Run the app

From the project root dir running `make fresh` will start the API in clean state. Any existing data & tables will be
dropped/recreated reverting the app to its out of the box state.

![Alt text](pic1.png?raw=true "Optional Title")

Running `make run` will start the API retaining the current database state.

### Querying the API

- Get a list of Area Profiles:
  ```shell
   curl -XGET "http://localhost:8080/profiles"
  ```

- Get an area profile by `area_code`:
  ```shell
  curl -XGET "http://localhost:8080/profiles/E05011362"
  ```

- Get a list of versions of an area profile (the default state has no previous versions - see next for details on 
  how to add one)
  ````shell
  curl -XGET "http://localhost:8080/profiles/E05011362/versions"
  ````

- Add a new version of key statistics to an area profile. When a new verison is added the "current"
  key stats values are copied into a version history table and then key stats table is updated with the latest values.  
  In this poc making a PUT request to this endpoint will reimport the same data again __i.e.__ all versions of the data 
  will be identical. In the real world these will be different and all or some values could be updated at different times 
  but here this fucntionality is intended to serve as an illustration onhow versioning the data can be achieved 
  using a version history table.
  ```shell
  curl -XPUT "http://localhost:8080/profiles/E05011362"
  ```
- Get a previous version of an area profiles by `version_id`:
  ````shell
  curl -XGET "http://localhost:8080/profiles/E05011362/versions/1000"
  ````

