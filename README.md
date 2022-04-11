# dp-area-profiles-design-spike

POC API demonstrating a proposed relational database schema to model Area Profiles & key statistics data. The POC allows
you to load 1 or more versions of test data which can then be queried via the API endpoints. 

:warning: This is a POC, not production ready code and is for illustrative purposes only. :warning:

### Getting started

Get the code:

  ```bash
  git clone https://github.com/ONSdigital/dp-area-profiles-design-spike
  ```

The POC is backed by a Postgres DB which needs to be set up before running the app. Ensure Docker is running then cd 
into the `v2` dir and run:

  ```
  make compose
  ```
Add the following env vars to your profile: 
````bash
export AP_POSTGRES_USER=postgres
export AP_POSTGRES_PASSWORD=mysecretpassword
export AP_DATABASE_NAME=area_profiles
````

Open another terminal and run the following to connect to Postgres:

```bash
docker exec -it v2_postgres_1 psql -U postgres
```

** :warning: The container name my vary. Use `docker ps` to get the name of your conatiner and replace
`v2_postgres_1` as required.

Create a database for the POC to connect to:

  ```
  CREATE DATABASE area_profiles;
  ```

### Run the app
`poc` is a simple _Cli_ with 2 commands:

- `init` - initalise / drop & recreate the area profiles database. For more details see the help command `./poc init -h`
- `api` - run the area profiles API.  For more details see the help command `./poc api -h`

Build the `poc` binary:
```bash
make build
```
Create and populate the database with 2 versions of test data for 1 area profile.
````bash
./poc init -l=1.csv -l=2.csv
````
Run the API (http://localhost:8080/profiles)
````bash
./poc api
````

### Querying the API

- **Get Area Profiles**:
  ```shell
  curl -XGET "http://localhost:8080/profiles"
  ...
  [
    {
      "id": 1000,
      "name": "Resident Population for Disbury East, Census 2021",
      "area_code": "E05011362",
      "href": "http://localhost:8080/profiles/E05011362"
    }
  ]
  ```

- **Get area profile** by `area_code`:
  ```shell
    curl -XGET "http://localhost:8080/profiles/E05011362"
    ...
    {
      "id": 1000,
      "name": "Resident Population for Disbury East, Census 2021",
      "area_code": "E05011362",
      "href": "http://localhost:8080/profiles/E05011362/stats"
    }
  ```

- **Get Area Profile (current) key stats** (some entries omitted)
  ````shell
    curl -XGET "http://localhost:8080/profiles/E05011362/stats"
    ...
    [
      {
        "id": 1100,
        "profile_id": 1000,
        "name": "Population density (Hectares)",
        "value": "1",
        "unit": "",
        "date_created": "2022-04-11T16:12:25.30247Z",
        "last_modified": "0001-01-01T00:00:00Z",
        "metadata": {
          "dataset_id": "efg789",
          "dataset_name": "Test dataset 2",
          "href": "http://localhost:8080/datasets/efg789"
        }
      },
      ...
      {
        "id": 1000,
        "profile_id": 1000,
        "name": "Resident population",
        "value": "2",
        "unit": "",
        "date_created": "2022-04-11T16:12:25.30247Z",
        "last_modified": "0001-01-01T00:00:00Z",
        "metadata": {
          "dataset_id": "abc123",
          "dataset_name": "Test dataset 1",
          "href": "http://localhost:8080/datasets/abc123"
        }
      }
    ]
  ````
- **Get Key Stats versions**
  ````shell
  curl -XGET "http://localhost:8080/profiles/E05011362/stats/versions"
  ...
  {
    "id": 1000,
    "name": "Resident Population for Disbury East, Census 2021",
    "area_code": "E05011362",
    "href": "http://localhost:8080/profiles/E05011362/stats",
    "versions": [
      "2022-04-11T16:12:25.332978Z",
      "2022-04-11T16:12:25.30247Z"
    ]
  }
  ````
- **Get Key stats by version** (some results omitted)
  ````shell
  curl -XGET "http://localhost:8080/profiles/E05011362/stats/versions/2022-04-11T16:12:25.30247Z"
  ...
  [
    {
      "id": 1100,
      "profile_id": 1000,
      "name": "Population density (Hectares)",
      "value": "1",
      "unit": "",
      "date_created": "2022-04-11T16:12:25.30247Z",
      "last_modified": "0001-01-01T00:00:00Z",
      "metadata": {
        "dataset_id": "efg789",
        "dataset_name": "Test dataset 2",
        "href": "http://localhost:8080/datasets/efg789"
      }
    },
    ...
    {
      "id": 1000,
      "profile_id": 1000,
      "name": "Resident population",
      "value": "1",
      "unit": "",
      "date_created": "2022-04-11T16:12:25.30247Z",
      "last_modified": "0001-01-01T00:00:00Z",
      "metadata": {
        "dataset_id": "abc123",
        "dataset_name": "Test dataset 1",
        "href": "http://localhost:8080/datasets/abc123"
      }
    }
  ]
  ````

