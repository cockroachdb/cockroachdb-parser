#!/bin/sh

DIR=$1
FROM=$2
TO=$3

mkdir -p $TO/$DIR
for file in $(find $FROM/$DIR -maxdepth 1 -type f); do
  cp $file $TO/$DIR/$(basename $file)
done
