package profiling

import (
	"cowboy-gorl/pkg/logging"
	"time"

    "net/http"
    _ "net/http/pprof"
)

// TimedCall executes the passed function, and prints the time it took to
// execute it using logging.Debug()
func TimedCall(f func(), function_name string) {
    before := time.Now()
    f()
    exec_time := time.Since(before)
    logging.Debug("Function %v took %v microsec to execute", function_name, exec_time.Microseconds())
}

// RunPProf launches an http server for PProf to hook into. Information can
// be fetched at `http://<adress>/debug/pprof`.
//
// Passing an empty string as address will default to localhost:8080
func RunPProf(address string) {
    if address == "" {
        address = "localhost:8080"
    }
    go func() {
        logging.Debug("Launching PProf Server")
        logging.Debug("%v", http.ListenAndServe(address, nil))
    }()
}
