# Remi

## Introduction

A code challenge from remi.

Try it at: http://remiproject.online/

## Usage

#### Prerequisite

golang and docker must be installed.

#### How to run the app
- Before run program, must run docker-compose file and run db migration files
    - Run docker-compose
  ```
  docker-compose up -d
  ```

- In order to run the app without building the binary file, please run following commands:

    - Run with default config
  ```
  go run main.go
  ```

- If you want to run the app by building binary file, please run following commands:

```
go build -o challenge
./challenge
```

#### How to test the app

- Access to golang directory and run command go test:

```
go test ./...
```

## Source code structure explanation

```
|-.github/ 
|-deployments/
|-features/
|-internal/
    |-entities/
    |-repositories/
    |-services/
|-migrations/
|-pkg/
    |-cmsql/
    |-config/
    |-crypto/
    |-golibs/
    |-xerror
|-templates/
|-up/
|-main.go 
```

The app is splitted into serveral packages:

- **.github**: It contains workflows of github action
    - workflow of unit test

- **migrations**: It contains migration files for database.

- **features**: It contains integration tests

- **internal**:
    - entities: Contains entities 
    - repositories: Interact with database
    - services: Provides functions to handle requests and return responses.

- **pkg**:
    - config: Provides functions to load config from file or default.
    - cmsql: Provides functions for config Postgres.
    - crypto: Hash and Check password.
    - golibs: Provide some common functions for database and go utils
    - xerror: Define errors and map its with httpStatus.

- **templates**: It contains frontend of this project

- **up**: Define structure responses and service interfaces.

- **main.go**: It load configs and run server.

## Future

- Create tool for log request, sql.
- Add dependencies injection.
- Add workflow to auto-deploy (github action)