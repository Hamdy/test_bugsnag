# test_bugsnag

A tiny Go API with a real bug, for trying pipe.ai's "fix the error and open a pull request" workflows end to end.

| Route | Answers |
|---|---|
| `GET /` | `200 {"status":"ok"}` |
| `GET /test` | `200 {"average_ms":0}`: the average of the recent latencies, 0 while none have been recorded. (It used to divide by zero on the empty list; any panic in a handler still goes to BugSnag with its stack trace and the request) |

## Run it

1. Copy `.env.example` to `.env` and set `BUGSNAG_API_KEY` (BugSnag → Project settings → Notifier API key).
   `.env` is ignored by git.
2. Start it, on port 8080 unless `PORT` says otherwise:

   ```bash
   make run
   ```

3. Try both routes:

   ```bash
   curl -s localhost:8080/
   ```

   ```bash
   curl -s localhost:8080/test
   ```

   The second answers `{"average_ms":0}`. A panic in a handler answers 500 and, with a key set, an event appears in
   BugSnag within a few seconds.

## Tests

```bash
make test
```

`TestAverageLatency` covers the happy path, `TestAverageLatencyEmpty` the empty list, and `TestLatencyReport` checks
that `GET /test` answers 200.
