# Dashboard KB with a bad dashboard path (fixture)

The Locations block declares a dashboard file that does not exist, so resolution
must fail with a clear "no dashboard file" error rather than falling through to
the planning-area check.

## Locations

dashboard: knowledge/does-not-exist.md
