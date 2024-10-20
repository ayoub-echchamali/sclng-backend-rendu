# Submission for "hard-skill" Scalingo


## Introduction


## Execution

```
docker compose up
```

Application will be then running on port `5000`

## Test

```
$ curl localhost:5000/ping
{ "status": "pong" }
```

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
