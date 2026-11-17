// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package pdf

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eclipse-disuko/disuko/logy"
)

var requestSessionTest = &logy.RequestSession{ReqID: "TEST"}

// isPDF reports whether buf starts with the PDF magic header.
func isPDF(buf []byte) bool {
	return len(buf) >= 4 && string(buf[:4]) == "%PDF"
}

// writeHTML writes a minimal HTML file with the given title into dir.
func writeHTML(t *testing.T, dir, name, title string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	content := "<html><body><h1>" + title + "</h1></body></html>"
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writing test html %s: %v", p, err)
	}
	return p
}

// skipIfNoChrome skips the test if ConvertAndMerge failed because no Chrome/
// Chromium binary is available in this environment (e.g. minimal CI images).
func skipIfNoChrome(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	msg := err.Error()
	if strings.Contains(msg, "exec:") || strings.Contains(msg, "executable file not found") ||
		strings.Contains(msg, "no such file or directory") {
		t.Skipf("skipping: Chrome/Chromium not available: %v", err)
	}
}

func TestConvertAndMergeSingleFile(t *testing.T) {
	dir := t.TempDir()
	htmlFile := writeHTML(t, dir, "page.html", "hello")
	out := filepath.Join(dir, "out.pdf")

	err := ConvertAndMerge(requestSessionTest, out, nil, nil, []string{htmlFile})
	skipIfNoChrome(t, err)
	if err != nil {
		t.Fatalf("ConvertAndMerge failed: %v", err)
	}

	buf, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output pdf missing: %v", err)
	}
	if len(buf) == 0 {
		t.Fatal("output pdf is empty")
	}
	if !isPDF(buf) {
		t.Fatalf("output is not a valid pdf, header: %q", buf[:min(4, len(buf))])
	}
}

func TestConvertAndMergeWithHeaderFooter(t *testing.T) {
	dir := t.TempDir()
	htmlFile := writeHTML(t, dir, "page.html", "hello")
	out := filepath.Join(dir, "out.pdf")
	header := "<span class=\"title\"></span>"
	footer := "<span class=\"pageNumber\"></span>"

	err := ConvertAndMerge(requestSessionTest, out, &header, &footer, []string{htmlFile})
	skipIfNoChrome(t, err)
	if err != nil {
		t.Fatalf("ConvertAndMerge failed: %v", err)
	}

	buf, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output pdf missing: %v", err)
	}
	if !isPDF(buf) {
		t.Fatalf("output is not a valid pdf, header: %q", buf[:min(4, len(buf))])
	}
}

func TestConvertAndMergeMultipleFiles(t *testing.T) {
	if _, err := exec.LookPath("pdfunite"); err != nil {
		t.Skip("skipping: pdfunite (poppler-utils) not available on PATH")
	}

	dir := t.TempDir()
	htmlFile1 := writeHTML(t, dir, "page1.html", "first")
	htmlFile2 := writeHTML(t, dir, "page2.html", "second")
	out := filepath.Join(dir, "out.pdf")

	err := ConvertAndMerge(requestSessionTest, out, nil, nil, []string{htmlFile1, htmlFile2})
	skipIfNoChrome(t, err)
	if err != nil {
		t.Fatalf("ConvertAndMerge failed: %v", err)
	}

	buf, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("merged output pdf missing: %v", err)
	}
	if !isPDF(buf) {
		t.Fatalf("merged output is not a valid pdf, header: %q", buf[:min(4, len(buf))])
	}
}
