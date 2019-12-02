package filetool

import (
	"bufio"
	"bytes"
	"testing"
)

func TestReadPagingBuffer(t *testing.T) {
	var (
		buf = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
		bs  []byte
		n   int
	)
	bs, n, _ = ReadPagingBuffer(0, 3, 0, bufio.NewReader(bytes.NewReader(buf)))
	if !bytes.Equal(buf[:3], bs[:n]) {
		t.Error()
	}
	bs, n, _ = ReadPagingBuffer(1, 3, 0, bufio.NewReader(bytes.NewReader(buf)))
	if !bytes.Equal(buf[3:6], bs[:n]) {
		t.Error()
	}
	bs, n, _ = ReadPagingBuffer(2, 3, 0, bufio.NewReader(bytes.NewReader(buf)))
	if !bytes.Equal(buf[6:9], bs[:n]) {
		t.Error()
	}
}
