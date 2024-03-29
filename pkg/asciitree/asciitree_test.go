package asciitree

import (
  "bytes"
  "testing"
  "github.com/stretchr/testify/assert"

  "htar/pkg/testdata"
  . "htar/pkg/core"
)

func TestNoPanic(t *testing.T) {
  part := testdata.SingleFilePart("test.txt", 42)
  buf := new(bytes.Buffer)

  tree := &PrintOptions{}
  tree.printParts(buf, 42, []Partition{ part })
  assert.True(t, buf.Len() > 0)
}

func TestNoDivideByZero(t *testing.T) {
  part := Partition{
    Groups: []FileGroup{
      FileGroup{},
    },
  }

  buf := new(bytes.Buffer)
  tree := &PrintOptions{FileCount: 42}
  tree.printParts(buf, 42, []Partition{ part })
  assert.True(t, buf.Len() > 0)
}
