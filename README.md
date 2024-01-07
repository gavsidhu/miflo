# Miflo

[![Release](https://img.shields.io/github/release/gavsidhu/miflo.svg)](https://github.com/gavsidhu/miflo/releases) [![Miflo CI/CD](https://github.com/gavsidhu/miflo/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/gavsidhu/miflo/actions/workflows/ci-cd.yml) [![Go Report](https://goreportcard.com/badge/github.com/gavsidhu/miflo)](https://goreportcard.com/report/github.com/gavsidhu/miflo)

Miflo is a database schema migration tool designed to simplify database schema changes for SQLite, PostgreSQL, and libSQL databases.

## Table of Contents
- [Installation](#installation)
- [Commands](#commands)
- [Usage](#usage)
  - [Connect Your Database](#connect-your-database)
    - [SQLite](#sqlite)
    - [PostgreSQL](#postgresql)
    - [libSQL](#libsql)
      - [Turso](#connecting-to-a-turso-database-instance)
      - [sqld](using-sqld)
  - [Create a migration](#Create-a-migration)
  - [Apply Migrations](#apply-migrations)
  - [Revert Migrations](#revert-migrations)
  - [List Migrations](#list-migrations)
- [Migrations](#migrations)
  - [Migration files](#migration-files)
  - [Migrations table](#migrations-table)
- [Contributing](#contributing)
- [License](#license)

## Installation

**Mac**

Install miflo on Mac by running the following command in your terminal:

```sh
curl -sSL https://github.com/gavsidhu/miflo/raw/main/scripts/install-mac.sh | bash
```

**Linux**

Install miflo on Mac by running the following command in your terminal:

```sh
curl -sSL https://github.com/gavsidhu/miflo/raw/main/scripts/install-linux.sh | bash
```

## Usage

### Connect Your Database

To connect your database using miflo, you need to set the DATABASE_URL in your .env file. This URL specifies the database type and its connection details. The format of the DATABASE_URL varies based on the type of database you are connecting to. Here's how you can set it up for each supported database:

#### SQLite

For SQLite, the DATABASE_URL is set using `sqlite:` followed by the path to your database file. Example:

```env
DATABASE_URL=sqlite:/path/to/your/databasefile.db
```

#### PostgreSQL

For connecting to a PostgreSQL database, the URL is formatted as `postgresql://username:password@host:port/database_name`. Replace the placeholders with your actual PostgreSQL database credentials. Example:

```env
DATABASE_URL=postgresql://user:password@localhost:5432/mydatabase
```

#### libSQL

For connecting to a libSQL database, miflo supports two protocols: HTTP and a custom libsql protocol. The appropriate protocol to use depends on your specific database setup.

##### Connecting to a Turso Database Instance

- If you are connecting to a Turso database instance, use the libsql protocol. The URL should be in the format `libsql://your_turso_database_url?authToken=your_auth_token`. Example:

```env
DATABASE_URL=libsql://your-db-instance.turso.io?authToken=yourActualAuthTokenHere
```

##### Using sqld

- When using sqld, the connection URL should use the HTTP protocol. The format will be something like `http://your_sqld_host:port`. Example:

```env
DATABASE_URL=http://127.0.0.1:8080
```

### Create a migration

Command: `miflo create [migration_file_name]`

- **Function**: The `create` command creates a new migration directory `[timestamp]_[migration_file_name]` in the migrations directory in the root of your project.
- **Migration Folder**: If a `migrations` folder doesn't exist at the root, miflo will prompt you to create one.
- **Genereated Files**: The command creates two SQL files in the new directory:
  - `up.sql` for applying the migration.
  - `down.sql` for reverting the migration.
- **Valid Migration Name**:
  - Migration names must adhere to a specific pattern. They should only contain letters (A-Z, a-z) and underscores (_).

Example:

```sh
miflo create add_users_table
```

### Apply migrations
Command: `miflo up`

- **Function** The `up` command applies all pending migrations. The migrations are applied in order of creation.
- **Execution**: miflo reads the `up.sql` files in each migration directory and executes the SQL statements contained within. These up.sql files define the changes to be made to the database schema.

```sh
miflo up
```

### Revert migrations
Command: `miflo revert`

- **Function**: The `revert` command is used to roll back applied migrations.
- **Revert Order**: Migrations are reverted in reverse order meaning the last applied migration is the first to be reverted.
- **Batch Operation**: This command uses batches. A batch is a group of migrations applied together during a single `miflo up` execution. `miflo revert` will only roll back the migrations of the last batch.
- **Execution**: miflo reads the `down.sql` files in each migration directory of the latest batch. These `down.sql` files contain SQL statements that undo the changes made by the corresponding `up.sql` files.

```sh
miflo revert
```

### List migrations
Command: `miflo list`

- **Function**: The `list` command lists all pending migrations.

```sh
miflo list
```

## Migrations 

### Migration Files

When a new migration is created using the miflo create command, it generates two SQL files in a dedicated directory for that migration. 

`up.sql`
- **Purpose**: The up.sql file is used for applying the migration. It contains SQL statements that modify the database schema, such as creating tables, adding columns, or other schema alterations.
- **Usage**: When you run the `miflo up` command, miflo executes the SQL statements in the `up.sql` files of each pending migration, in the order they were created. This process updates your database schema to the new desired state.
- **Content**: The content of an `up.sql` file typically includes `CREATE TABLE`, `ALTER TABLE`, `ADD COLUMN` statements, and other SQL commands that incrementally change the database structure.

`down.sql`
- **Purpose**: The down.sql file is used for reverting the migration. It contains SQL statements that undo the changes made by the up.sql file.
- **Usage**: When you run the `miflo revert` command, miflo executes the SQL statements in the `down.sql` files of the most recent batch of applied migrations, in reverse order. This process rolls back the latest changes made to your database schema.
- **Content**: The content of a `down.sql` file typically includes `DROP TABLE`, `DROP COLUMN`, and other SQL commands that reverse the changes made by the corresponding `up.sql`.

Example:
For a migration named `create_users_table`, the directory structure would be something like this:

```sh
/migrations
  /1704662056_create_users_table
    - up.sql
    - down.sql
```

### Migrations table

When you first use miflo to connect to your database, a table named `miflo_migrations` is automatically created. This table helps manage and track the state of database migrations.

**Table Columns**
- **id**: Serves as a unique identifier for each migration entry.
- **name**: Stores the name of the migration file. This is unique for each migration to prevent duplicate entries and to easily identify each migration.
- **batch**: Indicates the batch number in which the migration was applied. Migrations applied together in a single miflo up execution share the same batch number.
- **applied**: A boolean flag indicating whether the migration has been applied (true) or not (false).
- **applied_at**: Timestamp of when the migration was applied. It defaults to the current timestamp at the time of migration application.

## Contributing

If you want to contribute to miflo and make it better, your help is very welcome.

**Reporting Bugs**: If you encounter any bugs, please report them by opening a new issue.

**Suggesting Enhancements**: Got ideas to improve miflo? Feel free to open a new issue for discussion.

**Code Contributions**: You're welcome to fork the repository and submit pull requests.

**Testing with Databases**: To test miflo, you'll need to run Docker Compose, which will set up the necessary database environments for testing.

## License

Copyright 2023 Gavin Sidhu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
