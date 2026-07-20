// Command bklg serves a read-only kanban board for a knowledge-framework
// instance's planning and progress areas.
//
// This file is the CLI + server skeleton (TASK-001): argument splitting, flags,
// the loopback bind, and an exact-"/" handler that returns 200. Real area
// resolution, markdown parsing, and board/detail rendering arrive in later
// slices (see knowledge/reference/specs/bklg-spec.md §15).
package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gillchristian/bklg/internal/backlog"
)

// splitArgs pre-splits argv so the optional positional [path] works in any
// position while stdlib flag — which stops parsing at the first non-flag token
// — still sees every flag. Both "--port N" and "--port=N" forms are handled.
// Flags always land in flagArgs regardless of position; the first bare token is
// the path; any further bare tokens are returned in extra (the caller rejects
// them — the contract is a single optional [path]). Zero dependencies, spec §9.
func splitArgs(argv []string) (path string, flagArgs, extra []string) {
	path = "."
	takesValue := map[string]bool{"--port": true, "-port": true, "--dir": true, "-dir": true}
	seenPath := false
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		if strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
			if takesValue[a] && !strings.Contains(a, "=") && i+1 < len(argv) {
				i++ // the next token is this flag's value
				flagArgs = append(flagArgs, argv[i])
			}
			continue
		}
		if !seenPath {
			path, seenPath = a, true
		} else {
			extra = append(extra, a) // more than one [path] is malformed input
		}
	}
	return
}

func main() {
	path, flagArgs, extra := splitArgs(os.Args[1:])

	fs := flag.NewFlagSet("bklg", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), "usage: bklg [path] [--port N] [--dir D]\n")
		fs.PrintDefaults()
	}
	port := fs.Int("port", 1235, "port to listen on")
	dir := fs.String("dir", "knowledge", "knowledge dir, relative to [path]")
	if err := fs.Parse(flagArgs); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0) // -h/--help is a successful, requested action
		}
		os.Exit(2) // flag already printed the error + usage
	}
	if len(extra) > 0 {
		fmt.Fprintf(os.Stderr, "bklg: unexpected argument(s): %s (only one [path] is allowed)\n", strings.Join(extra, " "))
		fs.Usage()
		os.Exit(2)
	}

	areas, err := backlog.Resolve(path, *dir)
	if err != nil {
		// A root manifest (no planning area of its own, but a systems index) is
		// a helpful case, not a blank failure: list the per-system invocations.
		var rme *backlog.RootManifestError
		if errors.As(err, &rme) {
			fmt.Printf("bklg: no planning area at %s\n", rme.PlanningDir)
			fmt.Printf("This looks like a multi-system root manifest (%s).\n", rme.ManifestPath)
			fmt.Println("Point bklg at one system:")
			for _, s := range rme.Systems {
				fmt.Printf("  bklg %s --dir %s/knowledge\n", rme.Path, s)
			}
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "bklg: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, skeletonPage)
	})

	// Bind loopback only (spec §9) — this is a personal dev tool. Create the
	// listener before announcing readiness so "port in use" fails cleanly
	// without printing a false "Running" line.
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(*port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bklg: cannot listen on %s: %v\n", addr, err)
		os.Exit(1)
	}

	fmt.Printf("Running Backlog on port %d\n", *port)
	fmt.Printf("  knowledge: %s   planning: %s   progress: %s\n",
		backlog.DisplayPath(areas.KnowledgeDir), backlog.DisplayPath(areas.PlanningDir), backlog.DisplayPath(areas.ProgressDir))
	fmt.Printf("  http://localhost:%d\n", *port)

	if err := http.Serve(ln, mux); err != nil {
		fmt.Fprintf(os.Stderr, "bklg: server error: %v\n", err)
		os.Exit(1)
	}
}

const skeletonPage = `<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><title>bklg</title></head>
<body>
<h1>bklg</h1>
<p>Backlog board coming online — CLI + server skeleton (TASK-001).</p>
</body>
</html>
`
