package archive

import(
  "archive/tar"
  "bytes"
  "crypto/sha256"
  "errors"
  "fmt"
  "io"
  "strings"
  . "htar/pkg/core"
)

func VerifyPartition(reader io.Reader, progress chan<- ProgressUpdate) error {
  if progress != nil {
    defer close(progress)
  }

  tr := tar.NewReader(reader)

  hashes := make(map[string][]byte)
  var hashesBuf *bytes.Buffer

  meta := PartitionMeta{}
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

      switch header.Name {
      case partMetaFile:
        buf = new(bytes.Buffer)
      case hashesFile:
        hashesBuf = new(bytes.Buffer)
        buf = hashesBuf
      }

      read, hash, err := readFile(tr, buf)
      if err != nil {
        return fmt.Errorf("error reading file %q from archive: %v", header.Name, err)
      }

      if header.Name == partMetaFile {
        if err := meta.Decode(buf); err != nil {
          return fmt.Errorf("error parsing meta file: %v", err)
        }
      }

      if header.Name != hashesFile {
        hashes[header.Name] = hash
      }

      if buf == nil {
        // do not count meta files
        readFiles += 1
        readBytes += read
      }

      if progress != nil {
        pg := ProgressUpdate{
          Path: header.Name,
          Hash: hash,
          FileSize: read,
          CurrentFiles: readFiles,
          CurrentSize: readBytes,
          TotalFiles: meta.TotalFiles,
          TotalSize: meta.TotalSize,
        }
        if (buf != nil) {
          pg.CurrentFiles = 0
          pg.TotalSize = 0
        }
        progress <- pg
      }
    }
  }

  if hashesBuf == nil {
    return errors.New("archive does not contain checksums file")
  }

  expected, err := DecodeHashes(hashesBuf)
  if err != nil {
    return fmt.Errorf("error parsing checksum file: %v", err)
  }

  diff := CompareHashes(hashes, expected)

  if diff != nil {
    return fmt.Errorf("%d computed checksum did NOT match:\n%v", len(diff), strings.Join(diff, "\n"))
  }

  return nil
}

func readFile(tr *tar.Reader, content *bytes.Buffer) (int64, HexString, error) {
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
