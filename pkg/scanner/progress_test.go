package scanner

import (
  "bytes"
  "testing"
  "github.com/stretchr/testify/assert"
  "htar/pkg/testdata"
)

func TestInteractive(t *testing.T) {
  fsys := testdata.MakeTestFS()
  sources := []SourcePath{{Path: "var/pool", GroupingLevel: 2}}

  buf := new(bytes.Buffer)
  result, err := scanInteractive(fsys, buf, sources)

  assert.Nil(t, err)
  assert.Equal(t, 2, len(result))
  assert.True(t, buf.Len() > 0)
}
