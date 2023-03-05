package pipe

import (
  "bufio"
  "strings"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestSingleLine(t *testing.T) {
  scanner := bufio.NewScanner(strings.NewReader("Hello World\n"))
  scanner.Split(scanPtyLines)
  scanner.Scan()
  assert.Equal(t, "Hello World\n", scanner.Text())
}

func TestPartialLine(t *testing.T) {
  scanner := bufio.NewScanner(strings.NewReader("Hello World\rBye\n"))
  scanner.Split(scanPtyLines)
  scanner.Scan()
  assert.Equal(t, "Hello World", scanner.Text())
  scanner.Scan()
  assert.Equal(t, "\rBye\n", scanner.Text())
}

func TestStartCr(t *testing.T) {
  scanner := bufio.NewScanner(strings.NewReader("\rBye\rBye Bye\n"))
  scanner.Split(scanPtyLines)
  scanner.Scan()
  assert.Equal(t, "\rBye", scanner.Text())
  scanner.Scan()
  assert.Equal(t, "\rBye Bye\n", scanner.Text())
}

func TestStartLF(t *testing.T) {
  scanner := bufio.NewScanner(strings.NewReader("\n\nHello World\n"))
  scanner.Split(scanPtyLines)
  scanner.Scan()
  assert.Equal(t, "\n", scanner.Text())
  scanner.Scan()
  assert.Equal(t, "\n", scanner.Text())
  scanner.Scan()
  assert.Equal(t, "Hello World\n", scanner.Text())
}

func TestMixed(t *testing.T) {
  scanner := bufio.NewScanner(strings.NewReader("\rHey\nBye\rBye\n"))
  scanner.Split(scanPtyLines)
  scanner.Scan()
  assert.Equal(t, "\rHey\n", scanner.Text())
  scanner.Scan()
  assert.Equal(t, "Bye", scanner.Text())
  scanner.Scan()
  assert.Equal(t, "\rBye\n", scanner.Text())
}

