package archive

import(
  "bufio"
  "encoding/hex"
  "fmt"
  "io"
  "strings"
  "sort"
)

var hashesFile = "SHA256SUMS"

type Hashes map[string][]byte

func EncodeHashes(hashes Hashes, writer io.Writer) error {
  keys := make([]string, 0, len(hashes))
  for k := range hashes {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  for _, path := range keys {
    line := fmt.Sprintf("%x  %v\n", hashes[path], path);
    if _, err := io.WriteString(writer, line); err != nil {
      return err
    }
  }
  return nil
}

func DecodeHashes(reader io.Reader) (Hashes, error) {
  r := bufio.NewReader(reader)
  ln := 0
  hashes := make(Hashes)

  for true {
    line, _, err := r.ReadLine();
    ln += 1

    if err == io.EOF {
      break
    }
    
    if err != nil {
      return nil, err
    }

    if len(line) == 0 {
      // skip empty line
      continue
    }

    parts := strings.SplitN(string(line), " ", 2)

    if len(parts) != 2 || len(parts[1]) < 2 {
      return nil, fmt.Errorf("line %d improperly formatted", ln)
    }

    hash, err := hex.DecodeString(parts[0])
    if err != nil {
      return nil, fmt.Errorf("error decoding hash in line %d: %v", ln, err)
    }

    if prefix := parts[1][:1]; prefix != " " && prefix != "*" {
      return nil, fmt.Errorf("invalid prefix %q on line %d", prefix, ln)
    }

    // trim first character ' ' or '*' 
    name := parts[1][1:]

    hashes[name] = hash
  }

  return hashes, nil
}
