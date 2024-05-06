## Prerequisites

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [go-migrate](https://github.com/golang-migrate/migrate)

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/shamir92/project-sprint-team-happy.git
   cd project-sprint-team-happy/w1
   ```

2. Install dependencies

   ```bash
   go mod tidy
   ```

3. Set up the database:
   ```bash
   docker-compose up -d
   ```
4. Run database migrations
   ```bash
   make run-migrate-up
   ```

## Usage

1. Running the API

   ```bash
   go run .
   ```

2. Building the API

   ```bash
   make build/api
   ```

3. Creating a new migration file

   ```bash
   make run-migrate-create name=your_migration_name
   ```

4. Running database migrations

   To apply pending migrations, run:

   ```bash
   make run-migrate-up
   ```

   To rollback the last migration, run:

   ```bash
   make run-migrate-down-by-1
   ```

   To rollback all migrations, run:

   ```bash
   make run-migrate-down-all
   ```
