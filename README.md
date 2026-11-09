# test_bugsnag

A tiny Go API with a real bug, for trying pipe.ai's "fix the error and open a pull request" workflows end to end.

| Route | Answers |
|---|---|
| `GET /` | `200 {"status":"ok"}` |
| `GET /test` | `500`: it averages an empty list of latencies and divides by zero. The panic goes to BugSnag with its stack trace and the request |

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

   The second answers 500 and, with a key set, an event appears in BugSnag within a few seconds.

## Tests

```bash
make test
```

`TestAverageLatency` covers the happy path only. A fix for `/test` should make `averageLatency` handle an empty list
(and add a test for it) rather than hide the panic.
