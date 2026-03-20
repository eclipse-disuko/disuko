// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/helper"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/hash"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
)

type zipfile struct {
	os       *os.File
	osClosed bool
	w        *zip.Writer
	wClosed  bool
}

func initZipfile(path string) (*zipfile, error) {
	var (
		res zipfile
		err error
	)
	res.os, err = os.Create(path)
	if err != nil {
		return nil, err
	}
	res.w = zip.NewWriter(res.os)
	return &res, nil
}

func (z *zipfile) addDir(srcPath string) error {
	name := filepath.Base(srcPath)
	if err := z.createDir(name); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join(srcPath, e.Name()))
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		z.copy(filepath.Join(name, e.Name()), f)
		f.Close()
	}
	return nil
}

func (z *zipfile) createDir(path string) error {
	fh := zip.FileHeader{
		Name:     path + "/",
		Modified: time.Now(),
	}
	if _, err := z.w.CreateHeader(&fh); err != nil {
		return err
	}
	return nil
}

func (z *zipfile) writeBytes(path string, content []byte) error {
	fh := zip.FileHeader{
		Name:               path,
		Modified:           time.Now(),
		UncompressedSize64: uint64(len(content)),
	}
	fz, err := z.w.CreateHeader(&fh)
	if err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	if _, err = fz.Write(content); err != nil {
		return fmt.Errorf("writing content: %w", err)
	}
	return nil
}

func (z *zipfile) writeJson(path string, content any) (string, error) {
	fh := zip.FileHeader{
		Name:     path,
		Modified: time.Now(),
	}
	fz, err := z.w.CreateHeader(&fh)
	if err != nil {
		return "", fmt.Errorf("writing header: %w", err)
	}
	j, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling: %w", err)
	}
	if _, err := fz.Write(j); err != nil {
		return "", fmt.Errorf("writing json content: %w", err)
	}
	return hash.GetSha256Hash(j), nil
}

func (z *zipfile) executeTmpl(path string, t *template.Template, data any) error {
	fh := zip.FileHeader{
		Name:     path,
		Modified: time.Now(),
	}
	fz, err := z.w.CreateHeader(&fh)
	if err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	if err := t.Execute(fz, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

func (z *zipfile) copy(path string, reader io.Reader) (string, error) {
	fh := zip.FileHeader{
		Name:     path,
		Modified: time.Now(),
	}
	fz, err := z.w.CreateHeader(&fh)
	if err != nil {
		return "", fmt.Errorf("writing header: %w", err)
	}
	r := helper.NewSha256ReaderWrapper(reader)
	if _, err = io.Copy(fz, r); err != nil {
		return "", fmt.Errorf("copying: %w", err)
	}
	return r.GetHash(), nil
}

func (z *zipfile) close() {
	if !z.wClosed {
		if err := z.w.Close(); err != nil {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorClose), err.Error())
		}
		z.wClosed = true
	}
	if !z.osClosed {
		if err := z.os.Close(); err != nil {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorClose), err.Error())
		}
		z.osClosed = true
	}
}
