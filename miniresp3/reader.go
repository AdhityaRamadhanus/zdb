package miniresp3

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type Reader struct {
	br *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	br := bufio.NewReader(r)
	br.Reset(r)
	return &Reader{br: br}
}

func (r *Reader) ReadArrayHeader() (count int, err error) {
	line, err := r.br.ReadString('\n')
	if err != nil {
		return -1, err
	}

	if TypeRESP(line[0]) != RESPArray {
		return -1, errors.Wrap(err, "expected array header")
	}

	// escaped string equal to "*%d"
	scanStr := string(RESPArray) + "%d"
	if _, err := fmt.Sscanf(line, scanStr, &count); err != nil {
		return -1, errors.Wrap(err, "failed to parse array length")
	}

	return count, nil
}

func (r *Reader) ReadBulkString() (bulkStr string, err error) {
	line, err := r.br.ReadString('\n')
	if err != nil {
		return "", err
	}

	if TypeRESP(line[0]) != RESPBulkString {
		return "", errors.Wrap(err, "expected bulk string header")
	}

	var bulkStrSize int
	// escaped string equal to "$%d"
	scanStr := string(RESPBulkString) + "%d"
	if _, err := fmt.Sscanf(line, scanStr, &bulkStrSize); err != nil {
		return "", errors.Wrap(err, "failed to parse bulk string size")
	}

	bulkBytes, err := io.ReadAll(io.LimitReader(r.br, int64(bulkStrSize)))
	if err != nil {
		return "", errors.Wrap(err, "failed to read bulk string")
	}

	bulkStr = string(bulkBytes)
	if len(bulkStr) != bulkStrSize {
		return "", errors.New("Length of actual data not same as in meta data")
	}

	// Find the resp \r\n bytes
	if b, err := r.br.ReadByte(); err != nil || b != '\r' {
		return "", errors.New("Missing CR")
	}
	if b, err := r.br.ReadByte(); err != nil || b != '\n' {
		return "", errors.New("Missing LF")
	}

	return bulkStr, nil
}
