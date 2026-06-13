package wal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	ty "kvstore/internal/types"
	"log"
	"os"
)

type WAL struct {
	file   *os.File
	writer *bufio.Writer
}

type Record struct {
	Seq   uint32
	Op    ty.OpType
	Key   string
	Value string
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
// [seq | op | key_len | key | val_len | val | checksum]
func (w *WAL) Append(seq uint32, op ty.OpType, key, val string) error {
	buf := new(bytes.Buffer)

	// Start writing to buffer
	// write operation's sequence_no
	binary.Write(buf, binary.LittleEndian, seq)

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

	// Reads everything that's written inside buffer till now and
	// creates checksum value based on that
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

func (w *WAL) Recover() ([]*Record, error) {
	w.file.Seek(0, io.SeekStart)
	r := bufio.NewReader(w.file)

	var records []*Record

	for {
		var seq uint32
		if err := binary.Read(r, binary.LittleEndian, &seq); err != nil {
			break
		}

		var op ty.OpType
		if err := binary.Read(r, binary.LittleEndian, &op); err != nil {
			break
		}

		var key_len uint32
		if klErr := binary.Read(r, binary.LittleEndian, &key_len); klErr != nil {
			break
		}

		key := make([]byte, key_len)
		if kErr := binary.Read(r, binary.LittleEndian, &key); kErr != nil {
			break
		}

		var val_len uint32
		if vlErr := binary.Read(r, binary.LittleEndian, &val_len); vlErr != nil {
			break
		}

		val := make([]byte, val_len)
		if vErr := binary.Read(r, binary.LittleEndian, &val); vErr != nil {
			break
		}

		var stored_checksum uint32
		if cErr := binary.Read(r, binary.LittleEndian, &stored_checksum); cErr != nil {
			break
		}

		// Reconstructs checksum value based on seq|op|key_len|key|val_len|val
		// compares with existing one for corruption
		buf := new(bytes.Buffer)

		binary.Write(buf, binary.LittleEndian, seq)
		binary.Write(buf, binary.LittleEndian, op)
		binary.Write(buf, binary.LittleEndian, key_len)
		buf.Write(key)
		binary.Write(buf, binary.LittleEndian, val_len)
		buf.Write(val)

		calcChecksum := crc32.ChecksumIEEE(buf.Bytes())

		// compares currently calculated and stored checksum value
		if calcChecksum != stored_checksum {
			return nil, errors.New("corrupted WAL record detected")
		}

		records = append(records, &Record{
			Seq:   seq,
			Op:    op,
			Key:   string(key),
			Value: string(val),
		})
	}

	return records, nil
}
