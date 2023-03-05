package scanner

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestReadLevel0(t *testing.T) {
  fsys := makeTestFS()
  result, err := ReadLevel(fsys, "var", 0);
  assert.Nil(t, err)
  assert.Equal(t, 1, len(result))
  assert.Equal(t, "var", result[0])
}

func TestReadLevel1(t *testing.T) {
  fsys := makeTestFS()
  result, err := ReadLevel(fsys, "var/pool", 1);
  assert.Nil(t, err)
  assert.Equal(t, 1, len(result))
  assert.Equal(t, "var/pool/data", result[0])
}

func TestReadLevel3(t *testing.T) {
  fsys := makeTestFS()
  result, err := ReadLevel(fsys, "var", 3);
  assert.Nil(t, err)
  assert.Equal(t, 2, len(result))
  assert.Equal(t, "var/pool/data/Documents", result[0])
  assert.Equal(t, "var/pool/data/Music", result[1])
}
