package archive

import (
  "bytes"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
  a := Hashes{ "name1" : {0x12, 0x34, 0x56}, "name2" : {0x98, 0x76, 0x54} }

  buf := new(bytes.Buffer)
  err := EncodeHashes(a, buf)
  assert.Nil(t, err)

  b, err := DecodeHashes(buf)
  assert.Nil(t, err)
  assert.Equal(t, a, b)
}

func TestTrimPrefix(t *testing.T) {
  buf := new(bytes.Buffer)
  buf.WriteString("123456  a\n")
  buf.WriteString("987654 *b\n")

  hashes, err := DecodeHashes(buf)
  assert.Nil(t, err)

  expected := Hashes{ "a" : {0x12, 0x34, 0x56}, "b" : {0x98, 0x76, 0x54} }
  assert.Equal(t, expected, hashes)
}

func TestEmpty(t *testing.T) {
  buf := new(bytes.Buffer)
  buf.WriteString("\n")

  hashes, err := DecodeHashes(buf)
  assert.Nil(t, err)
  assert.Equal(t, 0, len(hashes))
}

func TestInvalid(t *testing.T) {
  buf := new(bytes.Buffer)
  buf.WriteString("123456  valid\n")
  buf.WriteString("invalid\n")

  hashes, err := DecodeHashes(buf)
  assert.Nil(t, hashes)
  assert.EqualError(t, err, "line 2 improperly formatted")
}

func TestInvalidPrefix(t *testing.T) {
  buf := new(bytes.Buffer)
  buf.WriteString("123456 xa\n")

  hashes, err := DecodeHashes(buf)
  assert.Nil(t, hashes)
  assert.EqualError(t, err, "invalid prefix \"x\" on line 1")
}

func TestInvalidHash(t *testing.T) {
  buf := new(bytes.Buffer)
  buf.WriteString("invalid  invalid\n")

  hashes, err := DecodeHashes(buf)
  assert.Nil(t, hashes)
  assert.Regexp(t, "^error decoding hash in line 1:.*$", err.Error())
}
