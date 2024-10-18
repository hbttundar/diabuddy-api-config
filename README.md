<p align="center"><img src="art/diabuddy.webp" alt="Diabuddy API Config Package"></p>

# Diabuddy API Config Package

## Overview
The `diabuddy-api-config` package is the central point for managing environment configurations across the Diabuddy ecosystem. This package utilizes the `EnvManager` internally to provide robust environment variable management for different APIs, such as `user_api`, `auth_api`, and `food_api`. It ensures consistent configuration handling across all APIs.

## Key Features
- Load environment variables from `.env` files based on the application's environment (`test`, `production`, `local`, etc.).
- Support for default values when environment variables are not set.
- Ability to cache environment variable values to reduce repeated lookups.
- Extendable defaults to allow for customization per API or package.
- Dynamic Database Connection String (DSN) support for multiple database types.

## Installation
To install the `diabuddy-api-config` package in your Go project:

```sh
$ go get github.com/hbttundar/diabuddy-api-config
```

## Usage

### Creating a New ApiConfig Instance
`ApiConfig` is the main entry point for managing configurations across different APIs. Here's how you can initialize `ApiConfig`:
```go
import (
    "github.com/hbttundar/diabuddy-api-config/config/apiconfig"
    "fmt"
)

// Creating a new ApiConfig instance
apiConfig, err := apiconfig.NewApiConfig()
if err != nil {
    panic(err) // handle the error appropriately
}
```

### Retrieving Environment Variables
To get the value of an App environment variable through `ApiConfig`:

```go
// Retrieves "APP_NAME" from the environment; falls back to default if necessary
appName := apiConfig.App.Get("APP_NAME")
fmt.Println("App Name:", appName)
```

You can also provide a default value at the call site, which will be used if the environment variable is not set:

```go
dbHost := apiConfig.DB.Get("DB_HOST", "localhost")
fmt.Println("Database Host:", dbHost)
```

### Using Cache
`ApiConfig` supports caching via the `EnvManager` to avoid repeated lookups:

```go
// Enable caching
apiConfig, _ := apiconfig.NewApiConfig(apiconfig.WithUseCache(true))
cachedValue := apiConfig.App.Get("CACHE_TEST_KEY")
fmt.Println("Cached Value:", cachedValue)
```

To clear the cache, use:

```go
apiConfig.App.ClearCache()
```

## Database Configuration Using DBConfig

The `diabuddy-api-config` package also allows you to easily generate database connection strings (DSNs) using the `DBConfig` for different popular databases like PostgreSQL, MySQL, SQL Server, and others. Instead of working directly with DSNs, developers can use `DBConfig` to simplify the setup process.

### Creating a DBConfig Instance
Here's how you can create a `DBConfig` instance for configuring a database:

```go
import (
    "github.com/hbttundar/diabuddy-api-config/config/dbconfig"
    "github.com/hbttundar/diabuddy-api-config/config/envmanager"
    "fmt"
)

func main() {
    // Initialize the EnvManager
    envManager, err := envmanager.NewEnvManager()
    if err != nil {
        fmt.Println("Error initializing EnvManager:", err)
        return
    }

    // Create a new DBConfig instance for MySQL
    dbConfig, err := dbconfig.NewDBConfig(envManager, dbconfig.WithType("mysql"))
    if err != nil {
        fmt.Println("Error creating DBConfig:", err)
        return
    }

    // Generate the connection string
    connString, err := dbConfig.ConnectionString()
    if err != nil {
        fmt.Println("Error generating connection string:", err)
        return
    }
    fmt.Println("MySQL Connection String:", connString)
}
```

### Specifying Additional DSN Parameters
You can also specify additional DSN parameters when creating a `DBConfig` instance:

```go
import (
    "github.com/hbttundar/diabuddy-api-config/config/dbconfig"
    "github.com/hbttundar/diabuddy-api-config/config/envmanager"
    "fmt"
)

func main() {
    // Initialize the EnvManager
    envManager, err := envmanager.NewEnvManager()
    if err != nil {
        fmt.Println("Error initializing EnvManager:", err)
        return
    }

    // Create a new DBConfig instance for PostgreSQL with additional parameters
    params := map[string]string{"sslmode": "disable", "timezone": "UTC"}
    dbConfig, err := dbconfig.NewDBConfig(envManager, dbconfig.WithType("postgres"), dbconfig.WithDsnParameters(params))
    if err != nil {
        fmt.Println("Error creating DBConfig:", err)
        return
    }

    // Generate the connection string
    connString, err := dbConfig.ConnectionString()
    if err != nil {
        fmt.Println("Error generating connection string:", err)
        return
    }
    fmt.Println("PostgreSQL Connection String:", connString)
}
```

### Supported Database Types
- **Postgres**
- **MySQL**
- **SQL Server**
- **Oracle**
- **MongoDB**
- **Redis**
- **Cassandra**

## Testing with ApiConfig
When testing, use the `.env.test` file to configure the test environment. You can create a `.env.test` file by copying from `.env.dist`:

```sh
$ cp .env.dist .env.test
```

The GitHub Action workflow also performs this step automatically to ensure the test environment is properly set up.

To run tests:

```sh
go test -v ./...
```

This will run all the tests, leveraging `ApiConfig` to load configurations and ensure your test environment is configured correctly.

## Configuration Options
- **WithEnvironment(string)**: Load a specific `.env` file based on the provided environment name, such as `test` or `production`.
- **WithUseDefault(bool)**: Whether to use default values if an environment variable is not set.
- **WithUseCache(bool)**: Enables caching to reduce lookup overhead.
- **WithConnectionStringOptions**: Dynamic generation of DSN for popular databases, allowing you to easily manage connections across PostgreSQL, MySQL, SQL Server, Oracle, MongoDB, Redis, and Cassandra.

## Example
```go
package main

import (
    "fmt"
    "github.com/hbttundar/diabuddy-api-config/config/apiconfig"
    "github.com/hbttundar/diabuddy-api-config/config/dbconfig"
    "github.com/hbttundar/diabuddy-api-config/config/envmanager"
)

func main() {
    // Initialize the ApiConfig
    apiConfig, err := apiconfig.NewApiConfig()
    if err != nil {
        fmt.Println("Error initializing ApiConfig:", err)
        return
    }

    // Get a value from the environment
    appEnv := apiConfig.App.Get("APP_ENV", "development")
    fmt.Println("App Environment:", appEnv)

    // Initialize the EnvManager
    envManager, err := envmanager.NewEnvManager()
    if err != nil {
        fmt.Println("Error initializing EnvManager:", err)
        return
    }

    // Create a DBConfig instance for MySQL
    dbConfig, err := dbconfig.NewDBConfig(envManager, dbconfig.WithType("mysql"))
    if err != nil {
        fmt.Println("Error creating DBConfig:", err)
        return
    }

    // Generate the connection string for MySQL
    connString, err := dbConfig.ConnectionString()
    if err != nil {
        fmt.Println("Error generating connection string:", err)
        return
    }
    fmt.Println("MySQL Connection String:", connString)
}
```

## License
This project is licensed under the MIT License. See the LICENSE file for details.

