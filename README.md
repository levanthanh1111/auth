# Backend Service

## Domain Driven Design Architecture:

```bash
.
├── application
│   ├── handler
│   │   ├── contracts
│   │   ├── users
│   │   └── withdraw_requests
│   ├── middleware
│   └── router
├── config
├── domain
│   ├── repository
│   │   └── users
│   └── usecase
│       ├── contracts
│       ├── users
│       └── withdraw_requests
├── external
│   ├── auth_server
│   ├── db
│   └── mailer
├── model
└── shared
    ├── base
    ├── context
    ├── log
    └── utils
```

- Application Layer:
    Serving HTTP request
    Validation Request
    Manage transaction consistency
    Prepare Data Model
- Domain Layer
    Usecase: Business logics implementation only
    Repository: Data Storage entity logics implementation
- External Adapter
    External service adapter: DB Connection, Mailer Interface Adapter, Auth Server Interface adapter,...
- Model
    General Application Data Model
- Shared
    Internal shared library and helper

## Configuation

The application will read all configuration from `config.yaml` file at the current running folder
For The configuration template, please refer `config.example.yaml`

The application will refer the config from environment variables. Therefore, if an configuration is setted by environment variable, that config in `config.yaml` file will not applied.

There are some rules to replace configuration by environment variables:

- The environment variables must have prefix as CONFIG_
- The environment variables all CAPITALIZED
- The configuration in `yaml` format parent.child:

    ```yaml
    parent:
        child: value
    ```
Will be transform to enviroment variable: CONFIG_PARENT_CHILD=value

- The configuration in `yaml` format list:

    ```yaml
    list:
        - element1
        - element2
    ```
Will be transform to enviroment variable: CONFIG_LIST=element1,element2

## Setting up database with Docker:
- Run the following command with Docker installed

```bash
$ docker run --name msf -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=auth -d postgres
```

## Run service:

- Setup the configuration as the instruction above

```bash
$ go build -o main .
$ ./main
```
- Or

```bash
$ go run main.go
```

- go to  http://localhost:9080/docs to see api document
