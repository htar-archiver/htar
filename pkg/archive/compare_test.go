package archive

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCompareHashes(t *testing.T) {
  a := Hashes{ "file1" : {0x00, 0x01, 0xff}, "file2" : {0x00, 0x02, 0xff} }
  b := Hashes{ "file1" : {0x00, 0x66, 0xff}, "file2" : {0x00, 0x02, 0xff} }
  c := Hashes{ "file1" : {0x00, 0x01, 0xff}, "fileX" : {0x00, 0x02, 0xff} }
  d := Hashes{ "file1" : {0x00, 0x01, 0xff}, "file2" : {0x00, 0x02, 0xff}, "file3" : {0x00, 0x66, 0xff} }

  diff := CompareHashes(a, a)
  assert.Nil(t, diff)

  diff = CompareHashes(a, b)
  assert.Equal(t, 1, len(diff))
  assert.Regexp(t, "^file1: checksum mismatch.*", diff[0])

  diff = CompareHashes(a, c)
  assert.Equal(t, 2, len(diff))
  assert.Regexp(t, "^file2: checksum mismatch.*", diff[0])
  assert.Regexp(t, "^fileX: unexpected file.*", diff[1])

  diff = CompareHashes(a, d)
  assert.Equal(t, 1, len(diff))
  assert.Regexp(t, "^file3: unexpected file.*", diff[0])
}
