// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
)

type HashedReadAllWrap struct {
	sourceReader io.ReadSeekCloser
	hasher       hash.Hash
}

func NewHashedReadAllWrap(sourceReader io.ReadSeekCloser) *HashedReadAllWrap {
	wrapper := new(HashedReadAllWrap)
	wrapper.sourceReader = sourceReader
	wrapper.hasher = sha256.New()
	return wrapper
}

func (w *HashedReadAllWrap) ReadAllAndRewind() ([]byte, error) {
	var (
		tmp []byte = make([]byte, 512)
		res []byte
		err error
		n   int
	)
	_, errSeek := w.sourceReader.Seek(0, io.SeekStart)
	if errSeek != nil {
		return nil, err
	}
	for ; err == nil; n, err = w.sourceReader.Read(tmp) {
		res = append(res, tmp[:n]...)
	}
	_, errSeek = w.sourceReader.Seek(0, io.SeekStart)
	if errSeek != nil {
		return nil, err
	}
	if err == io.EOF {
		return res, nil
	}
	return nil, err
}

func (w *HashedReadAllWrap) GetHash() (string, error) {
	content, err := w.ReadAllAndRewind()
	if err != nil {
		return "", err
	}
	w.hasher.Write(content)
	return hex.EncodeToString(w.hasher.Sum(nil)), nil
}

type Sha256ReaderWrapper struct {
	sourceReader io.Reader
	hasher       hash.Hash
}

func (wrapper *Sha256ReaderWrapper) Read(p []byte) (n int, err error) {
	bytesRead, err := wrapper.sourceReader.Read(p)
	if err == nil || err == io.EOF {
		wrapper.hasher.Write(p[:bytesRead])
	}
	return bytesRead, err
}

func (wrapper *Sha256ReaderWrapper) GetHash() string {
	return hex.EncodeToString(wrapper.hasher.Sum(nil))
}

func NewSha256ReaderWrapper(sourceReader io.Reader) *Sha256ReaderWrapper {
	wrapper := new(Sha256ReaderWrapper)
	wrapper.sourceReader = sourceReader
	wrapper.hasher = sha256.New()
	return wrapper
}
