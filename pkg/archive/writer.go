package archive

import(
  "archive/tar"
  "bytes"
  "crypto/sha256"
  "fmt"
  "io"
  "io/fs"
  . "htar/pkg/core"
)

func WritePartition(fsys fs.FS, part Partition, writer io.Writer, progress chan<- ProgressUpdate) error {
  if progress != nil {
    defer close(progress)
  }

  tw := tar.NewWriter(writer)
	defer tw.Close()

  hashes := make(Hashes, part.TotalFiles)
  totalSize := part.TotalSize
  writtenFiles := int(0)
  writtenBytes := int64(0)

  for _, group := range part.Groups {
    for _, file := range group.Files {
      written, hash, err := writeFile(tw, fsys, file.Path)
      if err != nil {
        return fmt.Errorf("error appending file %q to archive: %v", file.Path, err)
      }

      hashes[file.Path] = hash

      writtenFiles += 1
      writtenBytes += int64(written)

      changed := written - int64(file.Size)
      totalSize += changed
      
      if progress != nil {
        progress <- ProgressUpdate{
          Path: file.Path,
          Hash: hash,
          FileSize: file.Size,
          FileChangedSize: changed,
          CurrentFiles: writtenFiles,
          CurrentSize: writtenBytes,
          TotalFiles: part.TotalFiles,
          TotalSize: totalSize,
        }
      }
    }
  }

  buf := new(bytes.Buffer)
  if err := EncodeHashes(hashes, buf); err != nil {
    return fmt.Errorf("error creating checksum file %v: ", err)
  }

  if _, hash, err := writeBuffer(tw, buf, "SHA256SUMS"); err != nil {
    return fmt.Errorf("error writing checksum file %v: ", err)
  } else {
    if progress != nil {
      progress <- ProgressUpdate{
        Path: "SHA256SUMS",
        Hash: hash,
      }
    }
  }

  return nil
}

func writeFile(tw *tar.Writer, fsys fs.FS, path string) (int64, []byte, error) {
  fi, err := fs.Stat(fsys, path)
  if err != nil {
    return 0, nil, err
  }

  header, err := tar.FileInfoHeader(fi, "")
  if err != nil {
    return 0, nil, err
  }

  // Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory structure would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = path

  if err := tw.WriteHeader(header); err != nil {
    return 0, nil, err
  }

  f, err := fsys.Open(path)
  if err != nil {
    return 0, nil, err
  }

  defer f.Close()

  sha := sha256.New()
  mw := io.MultiWriter(sha, tw)

  written, err := io.Copy(mw, f)
  if err != nil {
    return 0, nil, err
  }

  return written, sha.Sum(nil), nil
}

func writeBuffer(tw *tar.Writer, content *bytes.Buffer, path string) (int64, []byte, error) {
  header := &tar.Header{
    Name: path,
    Mode: 0644,
    Size: int64(content.Len()),
    Typeflag: tar.TypeReg,
  }

  if err := tw.WriteHeader(header); err != nil {
    return 0, nil, err
  }

  sha := sha256.New()
  mw := io.MultiWriter(sha, tw)

  written, err := mw.Write(content.Bytes());
  if err != nil {
    return 0, nil, err
  }

  return int64(written), sha.Sum(nil), nil
}
