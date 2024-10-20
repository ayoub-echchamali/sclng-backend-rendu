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

From Github API Docs:

`
The primary rate limit for unauthenticated requests is 60 requests per hour.
`

To increase that number to 5000, please specify a github jwt token in the configuration file: `config.json`

