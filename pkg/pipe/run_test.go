package pipe

import (
  "io"
  "strings"
  "os/exec"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestRunTrue(t *testing.T) {
  cmd := exec.Command("/bin/true")
  err := Run(cmd, nil, make(chan string))
  assert.Nil(t, err)
}

func TestRunFalse(t *testing.T) {
  cmd := exec.Command("/bin/false")
  err := Run(cmd, nil, make(chan string))
  assert.NotNil(t, err)
}

func TestEcho(t *testing.T) {
  cmd := exec.Command("/bin/echo", "Hello World")
  out := make(chan string, 1)

  err := Run(cmd, nil, out)
  assert.Nil(t, err)
  assert.Equal(t, "Hello World\n", <- out)
}

func TestStdErr(t *testing.T) {
  cmd := exec.Command("/bin/sh", "-c", "echo Hello World >&2")
  out := make(chan string, 1)

  err := Run(cmd, nil, out)
  assert.Nil(t, err)
  assert.Equal(t, "Hello World\n", <- out)
}

func TestCat(t *testing.T) {
  cmd := exec.Command("/bin/cat")
  in := io.NopCloser(strings.NewReader("Hello\nWorld\n"))
  out := make(chan string, 2)

  err := Run(cmd, in, out)

  assert.Nil(t, err)
  assert.Equal(t, "Hello\n", <- out)
  assert.Equal(t, "World\n", <- out)
}
