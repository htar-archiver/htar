package archive

import(
  "io"
  "fmt"
)

type Hashes map[string][]byte

func (hashes *Hashes) Encode(writer io.Writer) error {
  for path, hash := range *hashes {
    line := fmt.Sprintf("%x  %v\n", hash, path);
    if _, err := io.WriteString(writer, line); err != nil {
      return err
    }
  }
  return nil
}
