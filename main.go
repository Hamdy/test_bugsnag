// Command test_bugsnag is a tiny API with a real bug, for trying pipe.ai's fix-the-error
// workflows end to end: GET / answers {"status":"ok"}, and GET /test crashes, which reports the
// crash to BugSnag (with its stack trace and the request) and answers 500.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bugsnag/bugsnag-go/v2"
	bugsnagerrors "github.com/bugsnag/bugsnag-go/v2/errors"
)

// Version is the build, sent to BugSnag as the app version.
var Version = "dev"

func main() {
	// Stack frames are sent relative to the repository's root (main.go, not /home/…/main.go),
	// so whatever reads the event — a person, or a workflow fixing it — finds the file.
	root, _ := os.Getwd()
	root += string(os.PathSeparator)
	bugsnag.Configure(bugsnag.Configuration{
		// The API key comes from $BUGSNAG_API_KEY (see .env.example). Without one, errors are
		// logged here and not sent.
		APIKey:          os.Getenv("BUGSNAG_API_KEY"),
		ReleaseStage:    envOr("BUGSNAG_RELEASE_STAGE", "development"),
		AppVersion:      Version,
		ProjectPackages: []string{"main", "github.com/Hamdy/test_bugsnag"},
		SourceRoot:      envOr("BUGSNAG_SOURCE_ROOT", root),
		// Panics in handlers are reported by recoverer; BugSnag's default panic handler would
		// re-run the program as a child process to watch it, which isn't needed here.
		PanicHandler: func() {},
	})
	if os.Getenv("BUGSNAG_API_KEY") == "" {
		log.Print("BUGSNAG_API_KEY isn't set: errors won't reach BugSnag")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", status)
	mux.HandleFunc("GET /test", latencyReport)

	addr := ":" + envOr("PORT", "8080")
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, recoverer(mux)))
}

// status says the API is up.
func status(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// latencyReport reports the average response time of the recent requests. There are none yet,
// so the average is 0.
func latencyReport(w http.ResponseWriter, _ *http.Request) {
	var samples []int // milliseconds; nothing has been recorded
	writeJSON(w, http.StatusOK, map[string]int{"average_ms": averageLatency(samples)})
}

// averageLatency is the mean of samples, in milliseconds, or 0 when there are none.
func averageLatency(samples []int) int {
	if len(samples) == 0 {
		return 0
	}
	total := 0
	for _, s := range samples {
		total += s
	}
	return total / len(samples)
}

// recoverer turns a panic in a handler into a BugSnag event and a 500, instead of the dropped
// connection Go's server gives a panic by default.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}
			// Made here, inside the recover, the error's stack trace still reaches the line
			// that panicked.
			err := bugsnagerrors.New(p, 2)
			ctx := bugsnag.AttachRequestData(r.Context(), r)
			if nerr := bugsnag.Notify(err, ctx, r, bugsnag.SeverityError,
				bugsnag.MetaData{"route": {"path": r.URL.Path}}); nerr != nil {
				log.Printf("bugsnag: %v", nerr)
			}
			log.Printf("panic on %s %s: %v", r.Method, r.URL.Path, p)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("internal error: %v", p),
			})
		}()
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
