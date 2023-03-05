package archive

import(
  "archive/tar"
  "bytes"
  "crypto/sha256"
  "io"
)

func VerifyPartition(reader io.Reader, progress chan<- ProgressUpdate) error {
  if progress != nil {
    defer close(progress)
  }

  tr := tar.NewReader(reader)

  hashes := make(map[string][]byte)
  hashesFile := new(bytes.Buffer)

  readFiles := int(0)
  readBytes := int64(0)

  for {
    header, err := tr.Next()

    if err == io.EOF {
      break
    }

    if err != nil {
      return err
    }

    if (header.Typeflag == tar.TypeReg) {
      var buf *bytes.Buffer
      if header.Name == "SHA256SUMS" {
        buf = hashesFile
      }

      read, hash, err := readFile(tr, buf)
      if err != nil {
        return err
      }

      hashes[header.Name] = hash

      readFiles += 1
      readBytes += read

      if progress != nil {
        progress <- ProgressUpdate{
          Path: header.Name,
          FileSize: read,
          CurrentFiles: readFiles,
          CurrentSize: readBytes,
        }
      }
    }
  }

  return nil
}

func readFile(tr *tar.Reader, content *bytes.Buffer) (int64, []byte, error) {
  sha := sha256.New()

  var mr io.Writer
  if content != nil {
    mr = io.MultiWriter(sha, content)
  } else {
    mr = sha
  }

  read, err := io.Copy(mr, tr)
  if err != nil {
    return 0, nil, err
  }

  return read, sha.Sum(nil), nil
}
