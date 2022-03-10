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

- **Get Area Profiles**:
  ```shell
   curl -XGET "http://localhost:8080/profiles"
  ```
  Response:
  ```json
  [
    {
      "profile_id": 1000,
      "name": "Resident Population for Disbury East, Census 2021",
      "area_code": "E05011362",
      "href": "http://localhost:8080/profiles/E05011362"
    }
  ]
  ```

- **Get area profile** by `area_code`:
  ```shell
  curl -XGET "http://localhost:8080/profiles/E05011362"
  ```
  
  Response:
  ```json
  {
    "id": 1000,
    "name": "Resident Population for Disbury East, Census 2021",
    "area_code": "E05011362",
    "key_stats": [
      {
        "id": 1200,
        "profile_id": 1000,
        "name": "Average (mean) age",
        "value": "29",
        "unit": "",
        "date_created": "2022-03-09T16:35:58.664157Z",
        "last_modified": "0001-01-01T00:00:00Z"
      },
      ...
    ]
  }
  ```

- **Add a new version** of key statistics to an area profile. When a new verison is added the "current"
  key stats values are copied into a version history table and then key stats table is updated with the latest values.  
  In this poc making a PUT request to this endpoint will reimport the same data again __i.e.__ all versions of the data
  will be identical. In the real world these will be different and all or some values could be updated at different times
  but here this fucntionality is intended to serve as an illustration onhow versioning the data can be achieved
  using a version history table.
  ```shell
  curl -XPUT "http://localhost:8080/profiles/E05011362"
  ```

- Get a list of versions of an area profile (the default state has no previous versions - you will need to add one 
  first - see previous step)
  ````shell
  curl -XGET "http://localhost:8080/profiles/E05011362/versions"
  ````
  Response:
  ```json
  [
    {
      "id": 1000,
      "profile_id": 1000,
      "version_id": 1000,
      "date_created": "2022-03-09T16:23:04.02231Z",
      "last_modified": "0001-01-01T00:00:00Z",
      "href": "http://localhost:8080/profile/E05011362/versions/1000"
    }
  ]
  ```


- **Get a version** by `version_id`:
  ````shell
  curl -XGET "http://localhost:8080/profiles/E05011362/versions/1000"
  ````
  Response:
  ```json
  [
    {
      "version_id": 1000,
      "id": 1000,
      "profile_id": 1000,
      "name": "Average (mean) age",
      "value": "29",
      "unit": "",
      "date_created": "2022-03-09T16:23:04.02231Z",
      "last_modified": "2022-03-09T16:23:44.52323Z"
    },
    ...
    ...
    ...
  ]
```

