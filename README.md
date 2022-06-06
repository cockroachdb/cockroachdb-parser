# cockroachdb-parser

`cockroachdb-parser` is a snapshot of the parser package and
all its dependencies from the [CockroachDB repo](repo). The
smaller package is Apache licensed and contains less dependencies
to pull in when configuring compared to `go get`.

The SHA this is based off is available in `version`.

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

Ensure the [repo](repo) is cloned in your $GOPATH, and then type:

```sh
./snapshot.sh
```

[repo]: https://github.com/cockroachdb/cockroach
