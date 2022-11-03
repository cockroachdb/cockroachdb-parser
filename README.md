# cockroachdb-parser

`cockroachdb-parser` is a snapshot of the parser package and
all its dependencies from the [CRDB repo][repo]. The
smaller package is Apache licensed and contains less dependencies
to pull in when configuring compared to `go get github.com/cockroachdb/cockroach`.

The SHA this is based off is available in `version`.

## Versioning

Versioning is done by CockroachDB version, with a `v0.` prepended.
For example, `v0.22.1.0.x.y` maps to `v22.1.0` in CockroachDB, where
`.x.y` maps to any subiterations.

## Custom patches

There are custom patches in `patches/` which gets applied to the repo.
This helps us customise the parser slightly for third party users.

## Example usage

```
import (
  ...
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
  ...
)

func Parse() error {
  ast, err := parser.ParseOne("SELECT 1")
  if err != nil {
    return err
  }
  // Do something with the AST
  _ = ast
}
```

## Generating a snapshot

Ensure the [CRDB repo][repo] is cloned in your $GOPATH, and then type:

```sh
./snapshot.sh
```

[repo]: https://github.com/cockroachdb/cockroach
