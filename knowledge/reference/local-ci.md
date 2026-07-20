# Local CI

The commands behind "run local CI" — verification gate 5 and the `pr` delivery
gates. All run from the repo root. **All must pass before a PR opens.**

## Run local CI (copy-paste)

```sh
go build ./...                     # compiles; exit 0
go vet ./...                       # static checks; exit 0
test -z "$(gofmt -l .)"            # formatting: gofmt -l lists unformatted files; must be empty
go test ./...                      # tests; exit 0 ("no test files" is fine until tests exist)
```

One-liner (stops at first failure, prints a marker):

```sh
go build ./... && go vet ./... && { u="$(gofmt -l .)"; [ -z "$u" ] || { echo "unformatted: $u"; false; }; } && go test ./... && echo "LOCAL-CI: PASS"
```

Interpretation:
- `gofmt -l .` prints the *paths* of files that are not gofmt-clean; empty
  output = clean. (Use `gofmt -w .` to fix.)
- `go test ./...` prints `?   pkg   [no test files]` for packages without tests
  and still exits 0 — expected for the earliest slices.

## Smoke test (behavior gate 3)

Build and drive the real binary; quote the output in the journal. Skeleton:

```sh
go build -o /tmp/bklg ./cmd/bklg
/tmp/bklg --port 9099 . &          # or point at testdata: /tmp/bklg testdata --port 9099
sleep 0.3
curl -s -o /dev/null -w '%{http_code}\n' http://127.0.0.1:9099/
kill %1
```

## Remote check (delivery gate D3)

_(none recorded — the repo has no remote CI/GitHub Actions in v1, so gate D3 is
vacuous. If a CI workflow is ever added (parking-lot item), record its
status/wait command here — e.g. `gh run watch --exit-status` — the first time it
exists; until then do not claim a green remote CI.)_
