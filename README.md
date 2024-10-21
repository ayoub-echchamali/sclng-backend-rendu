# Submission for "hard-skill" Scalingo


## Introduction
This is the submission that answers the Scanlingo hard skill test. The goal from this test is to create a REST API that can deliver the last 100 public repositories hosted on github.com.

# Project Contents

## api

API scaffolding specification

- [handlers.go](api/handlers.go)

Contains the controller functions that parse the request, executes the necessary functions to build the response, and returns the built object. Uses DTOs to build consistent response object.

- [routes.go](api/routes.go)

Defines the API routes.

- [server.go](api/server.go)

Creates server instance and listens with speicifed port inside the configuration.

## config

- [config.go](config/config.go)

Defines the structure of the configuration file, as well as the functions needed to read and parse the configuration files.

## dto

- [structs.go](dto/structs.go)

Defines the Data Transfer Objects (DTO) structures for the JSON response objects.

## githubapi

Contains the functions that poll and retrieve data from the github api. Contains also the structures to parse the objects received.

- [public_repos.go](githubapi/public_repos.go)
- [structs.go](githubapi/structs.go)


## util

Various auxiliary generic functions 

- [http.go](util/http.go)
- [util.go](util/util.go)

## Execution

```
docker compose up
```

Application will be then running on port `5000`

## Test
To test api is up:
```
$ curl localhost:5000/ping
{ "status": "pong" }
```

To test the `\publicGithubRepos`:
`
go test -v
`

## Implementation choices

### API Interactions
#### Rate Limiting Management:
    
- Introduced error handling for rate limit errors, returning appropriate messages instead of causing the application to crash.

#### Fetching Repositories:

- Created a function to fetch repositories from the GitHub API, including error handling to manage various API responses and ensure stability.
- Utilized a separate function for making HTTP requests, encapsulating the logic for setting headers and reading responses.

#### Concurrent Language Fetching:

- Implemented concurrent fetching of languages for each repository using goroutines, enhancing performance by allowing multiple API calls to occur simultaneously.

- Used a `sync.WaitGroup` to ensure that all goroutines complete before proceeding with further processing.

### Filtering Logic
#### Dynamic Filtering:
- Allowed filtering of repositories based on specified languages and licenses using flags (`filterByLanguage`, `filterByLicense`).

- Introduced separate function, `shouldIncludeRepo`, to determine whether to include a repository in the final results based on the filtering criteria.

#### Flexible Query Parameters:

- Implemented logic to handle multiple query parameters for languages, allowing users to filter repositories by one or more languages.

- Ensured that the filtering mechanism checks for matches against the specified criteria without excessive complexity.

### Concurrent Processing
#### Race Condition Handling:

- Implemented safeguards to prevent race conditions when accessing shared resources, ensuring data integrity during concurrent processing.

## Notes

From Github API Docs, these are the rate limiting rules:

`
The primary rate limit for unauthenticated requests is 60 requests per hour.
`

`
No more than 100 concurrent requests are allowed. This limit is shared across the REST API and GraphQL API.
`

`
 No more than 900 points per minute are allowed for REST API endpoints.
`

`
No more than 90 seconds of CPU time per 60 seconds of real time is allowed.
`

`
Make too many requests that consume excessive compute resources in a short period of time
`

These are the condition that may be relevant for this test. The polling rate is adjusted to try and respect these condition as well as possible.

To use the full features of this submissions, it is advised to create a personal token from Github.com to increase the default 60 request per hour limit rate to 5000. 

If you choose to do so, please create a configuration file: `config.json` by copying the provided example `config.example.json` and specify your token inside of it.
