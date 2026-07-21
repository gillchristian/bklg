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
	takesValue := map[string]bool{"--port": true, "-port": true, "--dir": true, "-dir": true, "--dashboard": true, "-dashboard": true}
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
		fmt.Fprint(fs.Output(), "usage: bklg [path] [--port N] [--dir D] [--dashboard FILE]\n")
		fs.PrintDefaults()
	}
	port := fs.Int("port", 1235, "port to listen on")
	dir := fs.String("dir", "knowledge", "knowledge dir, relative to [path]")
	dashboard := fs.String("dashboard", "", "single-file dashboard to read (relative to [path]); selects the dashboard adapter (ADR-0004)")
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

	var srv *backlog.Server
	var echo func() // prints the mode-appropriate resolution line

	if *dashboard != "" {
		// Explicit --dashboard flag: single-file dashboard adapter (ADR-0004),
		// no manifest/Locations lookup, no planning area required.
		areas, err := backlog.ResolveDashboard(path, *dir, *dashboard)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bklg: %v\n", err)
			os.Exit(1)
		}
		srv, echo = dashboardServer(areas)
	} else {
		areas, err := backlog.Resolve(path, *dir)
		switch {
		case err != nil:
			// A root manifest (no planning area of its own, but a systems index)
			// aggregates every system into one board instead of erroring (TASK-012).
			var rme *backlog.RootManifestError
			if !errors.As(err, &rme) {
				fmt.Fprintf(os.Stderr, "bklg: %v\n", err)
				os.Exit(1)
			}
			srv = backlog.NewMultiServer(rme.Path, rme.Systems)
			echo = func() {
				fmt.Printf("  aggregate: %d systems — %s\n", len(srv.Systems()), strings.Join(srv.Systems(), ", "))
			}
		case areas.DashboardFile != "":
			// A dashboard: key in the manifest selected the adapter (ADR-0004).
			srv, echo = dashboardServer(areas)
		default:
			srv = backlog.NewServer(areas)
			echo = func() {
				fmt.Printf("  knowledge: %s   planning: %s   progress: %s\n",
					backlog.DisplayPath(areas.KnowledgeDir), backlog.DisplayPath(areas.PlanningDir), backlog.DisplayPath(areas.ProgressDir))
			}
		}
	}

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
	echo()
	fmt.Printf("  http://localhost:%d\n", *port)

	if err := http.Serve(ln, srv.Routes()); err != nil {
		fmt.Fprintf(os.Stderr, "bklg: server error: %v\n", err)
		os.Exit(1)
	}
}

// dashboardServer builds a single-instance Server for a resolved dashboard-mode
// Areas and the startup echo line for it (ADR-0004).
func dashboardServer(a backlog.Areas) (*backlog.Server, func()) {
	srv := backlog.NewServer(a)
	echo := func() {
		fmt.Printf("  knowledge: %s   dashboard: %s\n",
			backlog.DisplayPath(a.KnowledgeDir), backlog.DisplayPath(a.DashboardFile))
	}
	return srv, echo
}
