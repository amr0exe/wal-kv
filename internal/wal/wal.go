package wal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"hash/crc32"
	ty "kvstore/internal/types"
	"log"
	"os"
)

type WAL struct {
	file   *os.File
	writer *bufio.Writer
}

func NewWAL() (*WAL, error) {
	f, err := os.OpenFile("log/wal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file:   f,
		writer: bufio.NewWriter(f),
	}, nil
}

func (w *WAL) Close() error {
	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.file.Close()
}

// Append appends SET|DEL actions to wal file
// where actions are appended in binary format
// format in which actions are stored is as follow:
// [op | key_len | key | val_len | val | checksum]
func (w *WAL) Append(op ty.OpType, key, val string) error {
	buf := new(bytes.Buffer)

	// Start writing to buffer
	// write operation_type
	binary.Write(buf, binary.LittleEndian, op)

	// write key_len
	key_len := len(key)
	binary.Write(buf, binary.LittleEndian, uint32(key_len))

	// write key
	buf.WriteString(key)

	// write val_len
	val_len := len(val)
	binary.Write(buf, binary.LittleEndian, uint32(val_len))

	// write val
	buf.WriteString(val)

	// write checksum
	c := crc32.ChecksumIEEE(buf.Bytes())
	binary.Write(buf, binary.LittleEndian, c)

	log.Printf("% x", buf.Bytes())

	_, err := w.writer.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.file.Sync()
}
