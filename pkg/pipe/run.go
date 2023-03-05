package pipe

import (
  "io"
  "os/exec"
  "sync"
  "syscall"
)

func Run(cmd *exec.Cmd, stdin io.Reader, outlines chan<- string) error {
  defer close(outlines)

  cmd.SysProcAttr = &syscall.SysProcAttr{
    // Detach controlling terminal
    Setsid: true,
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
  return cmd.Wait()
}
