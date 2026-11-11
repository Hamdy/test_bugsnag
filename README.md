# test_bugsnag

A tiny Go API for trying pipe.ai's "fix the error and open a pull request" workflows end to end.

| Route | Answers |
|---|---|
| `GET /` | `200 {"status":"ok"}` |
| `GET /test` | `200 {"average_ms":0}`: the average of the recent latencies, 0 while none are recorded. A panic in a handler goes to BugSnag with its stack trace and the request, and is answered with a 500 |

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

## Tests

```bash
make test
```

`TestAverageLatency` covers the happy path, `TestAverageLatencyEmpty` and `TestLatencyReport` cover an empty list
of latencies (which used to divide by zero and panic).
