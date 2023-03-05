package pipe

import (
  "bufio"
  "io"
  "os/exec"
  "sync"
  "syscall"
)

func Run(cmd *exec.Cmd, stdinBuf io.Reader, outlines chan<- string) error {
  defer close(outlines)

  cmd.SysProcAttr = &syscall.SysProcAttr{
    // Detach controlling terminal
    Setsid: true,
  }

  stdin, err := cmd.StdinPipe()
  if err != nil {
    return err
  }

  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return err
  }

  stderr, err := cmd.StderrPipe()
  if err != nil {
    return err
  }

  var wg sync.WaitGroup

  wg.Add(1)
  go func() {
    defer wg.Done()
    readPipe(stdout, outlines)
  }()

  wg.Add(1)
  go func() {
    defer wg.Done()
    readPipe(stderr, outlines)
  }()

  if err := cmd.Start(); err != nil {
    return err
  }

  defer wg.Wait()

  if (stdinBuf != nil) {
    if _, err := io.Copy(stdin, stdinBuf); err != nil {
      return err
    }
  }

  stdin.Close()
  return cmd.Wait()
}

func readPipe(r io.Reader, lines chan<- string) {
  scanner := bufio.NewScanner(r)
  scanner.Split(scanPtyLines)
  for scanner.Scan() {
    lines <- scanner.Text()
  }
}
