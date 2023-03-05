package archive

import(
  "archive/tar"
  "bytes"
  "crypto/sha256"
  "errors"
  "fmt"
  "io"
  "strings"
)

func VerifyPartition(reader io.Reader, progress chan<- ProgressUpdate) error {
  if progress != nil {
    defer close(progress)
  }

  tr := tar.NewReader(reader)

  hashes := make(map[string][]byte)
  var hashesFile *bytes.Buffer

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
        hashesFile = new(bytes.Buffer)
        buf = hashesFile
      } else if header.Name == ".htar" {
        buf = new(bytes.Buffer)
      }

      read, hash, err := readFile(tr, buf)
      if err != nil {
        return fmt.Errorf("error reading file %q from archive: %v", header.Name, err)
      }

      if header.Name != "SHA256SUMS" {
        hashes[header.Name] = hash
      }

      if buf == nil {
        // do not count meta files
        readFiles += 1
        readBytes += read
      }

      if progress != nil {
        progress <- ProgressUpdate{
          Path: header.Name,
          Hash: hash,
          FileSize: read,
          CurrentFiles: readFiles,
          CurrentSize: readBytes,
        }
      }
    }
  }

  if hashesFile == nil {
    return errors.New("archive does not contain checksums file")
  }

  expected, err := DecodeHashes(hashesFile)
  if err != nil {
    return fmt.Errorf("error parsing checksum file: %v", err)
  }

  diff := CompareHashes(hashes, expected)

  if diff != nil {
    return fmt.Errorf("%d computed checksum did NOT match:\n%v", len(diff), strings.Join(diff, "\n"))
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
