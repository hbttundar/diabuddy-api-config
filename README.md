<p align="center"><img src="art/diabuddy.webp" alt="Diabuddy Error package"></p>
# ApiConfig Package

## Overview
The `diabuddyApiConfig` package is the central point for managing environment configurations across the Diabuddy ecosystem. This package utilizes the `EnvManager` internally to provide robust environment variable management for different APIs, such as `user_api`, `auth_api`, and `food_api`. It ensures consistent configuration handling across all APIs.

## Key Features
- Load environment variables from `.env` files based on the application's environment ( `test` or `production` or `local`, etc.).
- Support for default values when environment variables are not set.
- Ability to cache environment variable values to reduce repeated lookups anyhow, probably in the next version this will be removed. 
- Extendable defaults to allow for customization per API or package.

## Installation
To install the `ApiConfig` package in your Go project:
```sh
$ go get github.com/hbttundar/diabuddy-user-api
```

## Usage

### Creating a New ApiConfig Instance
`ApiConfig` is the main entry point for managing configurations across different APIs. Here's how you can initialize `ApiConfig`:
```go
import (
    "github.com/hbttundar/diabuddy-user-api/config/api"
    "fmt"
)

// Creating a new ApiConfig instance
apiConfig, err := api.NewApiConfig()
if err != nil {
    panic(err) // handle the error appropriately
}
```

### Retrieving Environment Variables
To get the value of an App environment variable through  `ApiConfig`:
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
apiConfig, _ := api.NewApiConfig(api.WithUseCache(true))
cachedValue := apiConfig.App.Get("CACHE_TEST_KEY")
fmt.Println("Cached Value:", cachedValue)
```

To clear the cache, use:
```go
apiConfig.App.ClearCache()
```

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

## Example
```go
package main

import (
    "fmt"
    "github.com/hbttundar/diabuddy-user-api/config/api"
)

func main() {
    // Initialize the ApiConfig
    apiConfig, err := api.NewApiConfig()
    if err != nil {
        fmt.Println("Error initializing ApiConfig:", err)
        return
    }

    // Get a value from the environment
    appEnv := apiConfig.App.Get("APP_ENV", "development")
    fmt.Println("App Environment:", appEnv)
}
```

## Contributing
Feel free to contribute to this repository. Fork it, make your changes, and create a pull request. Make sure your code follows the Go best practices and is properly tested.

## License
This project is licensed under the MIT License. See the LICENSE file for more information.

## Support
If you have any questions or need help, please feel free to create an issue on the GitHub repository or contact the maintainers.