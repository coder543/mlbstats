# mlbstats

## Overview

This is a server (called `mlbstats`) that will proxy requests for `/api/v1/schedule` to the mlb `statsapi`, with some modifications to both the request and response. Requests sent to `mlbstats` should be sent to the `/api/v1/schedule` route with two query parameters: `date` and `teamId`. In response, it will provide the usual response from `statsapi`, modified to place games from your team at the top of the list, including logic for handling various kinds of double headers.

`date` should be a date string formatted as `YYYY-MM-DD`. `teamId` should be the ID of a team, as listed in the `statsapi` on the `/api/v1/teams` route. That route is not proxied by `mlbstats` at this time.

## Building and running

```bash
go build
./mlbstats
```

This will start the proxy server listening on port 8080.

In a separate terminal, an example of a valid request would look like this:

```bash
curl 'http://localhost:8080/api/v1/schedule?date=2022-10-04&teamId=115'
```

`mlbstats` takes a couple of optional parameters:

```
Usage of mlbstats:
  -addr string
        The address used for the listening socket (default ":8080")
  -mock
        Controls whether this instance will use mocked data or real data
```

If you supply `-mock`, then it will read `schedule.json` from the current working directory and use that instead of proxying the `statsapi`.

## Tests

To run the built-in tests, simply run `go test ./...` from the root of the repo.

## TODO

`GameStatusInProgress` is currently just using a value that I selected, but it would be good to figure out what the MLB `statsapi` actually returns for the `StatusCode` when the game is in progress, and update this constant.
