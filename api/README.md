# Econest API

This is the backend API for **Econest**, an e-commerce platform. It is built with Go and structured as part of a monorepo architecture. This service handles core backend functionality such as product management, authentication, user services, and more.

---

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)

  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Super admin CLI](#super-admin-cli)
  - [Running the Server](#running-the-server)

- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- RESTful API built with Go
- User authentication and authorization
- Product and order management
- PostgreSQL database integration
- Key rotation for authentication token key pairs
- Scalable configuration system

---

## Getting Started

### Prerequisites

- Go (>= 1.23)
- PostgreSQL
- Redis
- `make` (optional, for using provided Makefile scripts)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/SaeedAlian/econest.git
   cd econest/api
   ```

2. Copy and configure the environment variables:

   ```bash
   cp .env.example .env
   # Then edit .env to match your local setup
   ```

3. Install Go dependencies:

   ```bash
   go mod tidy
   ```

### Super admin CLI

This API works with a hierarchy of user roles & permissions. This assures integrity and security
through the API routes. There is a special role called 'Super Admin' which can access all the
database and the super admin cli without any issues, but this user cannot log into the site via api,
and cannot be registered via the api, for security reasons. To create a super admin and an admin
user (for logging into the website with an admin account) please use the super admin cli first.
You can run the cli with this command:

```bash
go build -o bin/econestapi main.go && ./bin/econestapi --cli
```

Or with Makefile:

```bash
make run-super-admin-cli
```

### Running the Server

```bash
go run main.go
```

Or using the Makefile:

```bash
make run
```

Or

```bash
make run-prod
```

---

## Environment Variables

The `.env` file is used to configure the following variables:

- `PORT` - Port on which the server will run
- `DATABASE_URL` - PostgreSQL connection URL
- `JWT_SECRET` - Secret for signing JWT tokens
- `ENV` - Application environment (development, production, etc.)

Refer to `.env.example` for the full list of variables.
You can define ENV variable at the start to determine which env file you want to use.
For example:

```bash
ENV="devel" ./bin/econestapi
```

For using devel environment

---

## Project Structure

```
api/
├── main.go            # Entry point
├── api/               # API routes & subroutes declarations
├── config/            # Configuration management
├── db/                # Database migrations and DB manager
├── docs/              # Swagger documentations
├── lib/               # Shared internal libraries (such as dynamic menu)
├── services/          # Core services logic (auth, products, etc.)
├── types/             # Type definitions
├── utils/             # Utility functions
├── .env.example       # Sample environment variables
├── Makefile           # Helpful make commands
├── LICENSE            # Project license
├── README.md          # This file
├── .gitignore         # gitignore file
├── go.mod             # Module definitions & requirements
├── go.sum             # Dependency checksums
```

---

## Available Scripts

From the `api` directory, you can use:

- `make run` - Run the server in development mode
- `make run-prod` - Run the server in production mode
- `make build` - Build the project
- `make test` - Run the tests
- `make run-super-admin-cli` - Run the super admin CLI
- `make migration {migration_name}` - Generate a new migration file pair
- `make migrate-up` - Run the migrations
- `make migrate-down` - Rollback the migrations

---

## API Documentation

Interactive API docs are available via Swagger UI.

- Swagger UI: http://localhost:5000/swagger/index.html
- Swagger JSON spec: http://localhost:5000/swagger/doc.json

**NOTE**: You can replace '5000' with your defined port in the config variables

Documentation is generated using swaggo/swag. To regenerate docs after updating annotations:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

---

## Contributing

Feel free to submit issues or pull requests. Make sure to follow the coding conventions and test before pushing.

---

## License

This project is licensed under the GNU GPL V3.0. See the [LICENSE](../LICENSE) file for more information.
