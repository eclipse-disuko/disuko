// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package pdf

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/eclipse-disuko/disuko/logy"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

var requestSessionTest = &logy.RequestSession{ReqID: "TEST"}

// isPDF reports whether buf starts with the PDF magic header.
func isPDF(buf []byte) bool {
	return len(buf) >= 4 && string(buf[:4]) == "%PDF"
}

// hasChrome reports whether a Chrome/Chromium binary is available, mirroring
// the candidate names and paths chromedp checks in its execLocations.
func hasChrome() bool {
	var candidates []string
	switch runtime.GOOS {
	case "darwin":
		candidates = []string{
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
	case "windows":
		candidates = []string{
			"chrome",
			"chrome.exe",
			"msedge",
			"msedge.exe",
		}
	default:
		candidates = []string{
			"headless_shell",
			"headless-shell",
			"chromium",
			"chromium-browser",
			"google-chrome",
			"google-chrome-stable",
			"google-chrome-beta",
			"google-chrome-unstable",
			"/usr/bin/google-chrome",
			"/usr/local/bin/chrome",
			"/snap/bin/chromium",
			"chrome",
		}
	}
	for _, c := range candidates {
		if _, err := exec.LookPath(c); err == nil {
			return true
		}
	}
	return false
}

// skipIfNoChrome skips the test if no Chrome/Chromium binary is available in
// this environment (e.g. minimal CI images).
func skipIfNoChrome(t *testing.T) {
	t.Helper()
	if !hasChrome() {
		t.Skip("skipping: Chrome/Chromium not available")
	}
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

// assertValidPDF checks that path exists and contains a valid PDF document.
func assertValidPDF(t *testing.T, path string) {
	t.Helper()
	buf, err := os.ReadFile(path)
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

// assertPageCount checks that the PDF at path has the expected page count.
func assertPageCount(t *testing.T, path string, want int) {
	t.Helper()
	pages, err := api.PageCountFile(context.Background(), path)
	if err != nil {
		t.Fatalf("counting pages of %s failed: %v", path, err)
	}
	if pages != want {
		t.Fatalf("expected %d pdf pages in %s, got %d", want, path, pages)
	}
}

func TestConvertAndMergeSingleFile(t *testing.T) {
	skipIfNoChrome(t)

	dir := t.TempDir()
	htmlFile := writeHTML(t, dir, "page.html", "hello")
	out := filepath.Join(dir, "out.pdf")

	if err := ConvertAndMerge(requestSessionTest, out, nil, nil, []string{htmlFile}); err != nil {
		t.Fatalf("ConvertAndMerge failed: %v", err)
	}

	assertValidPDF(t, out)
	assertPageCount(t, out, 1)
}

func TestConvertAndMergeWithHeaderFooter(t *testing.T) {
	skipIfNoChrome(t)

	dir := t.TempDir()
	htmlFile := writeHTML(t, dir, "page.html", "hello")
	out := filepath.Join(dir, "out.pdf")
	header := "<span class=\"title\"></span>"
	footer := "<span class=\"pageNumber\"></span>"

	if err := ConvertAndMerge(requestSessionTest, out, &header, &footer, []string{htmlFile}); err != nil {
		t.Fatalf("ConvertAndMerge failed: %v", err)
	}

	assertValidPDF(t, out)
}

func TestConvertAndMergeMultipleFiles(t *testing.T) {
	skipIfNoChrome(t)
	if _, err := exec.LookPath("pdfunite"); err != nil {
		t.Skip("skipping: pdfunite (poppler-utils) not available on PATH")
	}

	dir := t.TempDir()
	htmlFile1 := writeHTML(t, dir, "page1.html", "first")
	htmlFile2 := writeHTML(t, dir, "page2.html", "second")
	out := filepath.Join(dir, "out.pdf")

	if err := ConvertAndMerge(requestSessionTest, out, nil, nil, []string{htmlFile1, htmlFile2}); err != nil {
		t.Fatalf("ConvertAndMerge failed: %v", err)
	}

	assertValidPDF(t, out)
	assertPageCount(t, out, 2)
}
