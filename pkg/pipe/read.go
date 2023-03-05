package pipe

import (
  "bytes"
)

func scanPtyLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
  if i := bytes.IndexByte(data, '\r'); i >= 0 {
    if (i >= 1) {
      // We have a partial line beeing updated
      return i, data[:i], nil
    }
    // Split before next CR
    cr := bytes.IndexByte(data[1:], '\r')
    lf := bytes.IndexByte(data[1:], '\n')
    if cr >= 0 && lf > cr {
      return cr + 1, data[:cr + 1], nil
    }
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
    // We have a full newline-terminated line.
    return i + 1, data[:i + 1], nil
	}
  // Return partial line
  return len(data), data, nil
}
