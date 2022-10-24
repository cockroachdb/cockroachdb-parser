#!/bin/sh

if [ -z ${GOPATH} ]; then
  echo "$GOPATH must be set"
  exit 1
fi

LOC=$(pwd)
COCKROACHDB_LOC=$GOPATH/src/github.com/cockroachdb/cockroach

echo "Generating files in-line using dev build short"
cd $COCKROACHDB_LOC
./dev build short

echo "Copying required dependencies over"
rm -rf $LOC/pkg
bazel query 'deps(//pkg/sql/parser)' | \
  # only care about dependencies in the CockroachDB repo
  grep '^//pkg' | \
  # remove the // prefix
  sed -e 's_^//__' | \
  # remove the :build_target directive
  sed -e 's/:.*$//' | \
  # remove testutils
  grep -v 'pkg/testutils' | grep -v 'pkg/cmd/protoc-gen-gogoroach' | \
  # ensure everything is unique
  sort | uniq | \
  # copy over each dependency
  xargs -I from_dir $LOC/copy_files.sh from_dir $COCKROACHDB_LOC $LOC
echo $(git rev-parse HEAD) > $LOC/version

echo "removing excess files, fixing up file paths"
cd $LOC
# temporarily copy embedded data over
cp -R $COCKROACHDB_LOC/pkg/geo/geoprojbase/data $LOC/pkg/geo/geoprojbase/data
# delete gitignores
rm pkg/sql/lexbase/.gitignore pkg/sql/parser/.gitignore
# delete tests and testfiles
find pkg -type f -name '*_test.go' | xargs rm
# sed replace any instances of cockroachdb
find pkg -type f -name '*.go' | xargs sed -i '' -e 's_github\.com/cockroachdb/cockroach/_github.com/cockroachdb/cockroachdb-parser/_g'

echo "Cleaning up go and testing everything works"
go mod tidy
git clean -fd
go test .
