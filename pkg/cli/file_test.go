package cli

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
  archiver := &FileArchiver{ Destination: "my/dir/test.tar" }
  assert.Equal(t, "my/dir/test_part42.tar" , archiver.getName(42))
}
