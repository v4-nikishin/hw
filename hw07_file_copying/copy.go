package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFile              = errors.New("the same file")
)

func Copy(from, to string, offset, limit int64) (err error) {
	sfi, err := os.Stat(from)
	if err != nil {
		return
	}
	if offset > sfi.Size() {
		return ErrOffsetExceedsFileSize
	}
	if !sfi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	dfi, err := os.Stat(to)
	if os.SameFile(sfi, dfi) {
		return ErrSameFile
	}
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return ErrUnsupportedFile
		}
	}
	err = copyContent(from, to, offset, limit, sfi.Size())
	return
}

func copyContent(src, dst string, offset, limit int64, fileSize int64) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	var totalN int64
	lim := limit
	if lim == 0 {
		lim = fileSize
	}
	if lim > fileSize-offset {
		lim = fileSize - offset
	}

	// create and start new bar
	bar := pb.Start64(lim)
	defer bar.Finish()

	bufSize := 10
	if int64(bufSize) > lim {
		bufSize = int(lim)
	}
	off := offset
	buf := make([]byte, bufSize)
	for {
		n, err := in.ReadAt(buf, off)
		if err != nil {
			if err == io.EOF {
				out.Write(buf[:n])
				bar.Add(n)
				break
			}
		}
		out.Write(buf[:n])
		bar.Add(n)
		totalN += int64(n)
		if limit != 0 && (totalN == lim) {
			break
		}
		off += int64(n)
	}
	return
}
