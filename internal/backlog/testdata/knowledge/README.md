# knowledge/ — demo instance (bklg fixture)

A small, complete framework instance used both as parser test input and as the
live demo bklg renders: `bklg internal/backlog/testdata`.

It has no `## Locations` block, so resolution falls back to `knowledge/planning`
and `knowledge/progress`. It deliberately seeds the three reconciliation
warnings and every badge case — see the planning/progress files.
