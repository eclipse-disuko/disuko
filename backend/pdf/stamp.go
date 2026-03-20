// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package pdf

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/temp"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func StampPageNumbers(rs *logy.RequestSession, path string) error {
	pageCount, err := api.PageCountFile(path)
	if err != nil {
		return fmt.Errorf("getting page count: %w", err)
	}

	fontFile := conf.Config.Server.GetDocumentResourceFilePath(conf.Config.Fonts.Stamp)
	if err = api.InstallFonts([]string{fontFile}); err != nil {
		return fmt.Errorf("installing font: %w", err)
	}

	watermarks := make(map[int]*model.Watermark)
	for i := range pageCount {
		indicator := fmt.Sprintf("%d/%d", i+1, pageCount)
		w, err := api.TextWatermark(indicator, "pos: bl, rot:0, points: 8, scale: 0.02 rel, offset: 30 30, color: 0 0 0, fontname: Roboto-Regular", false, false, types.POINTS)
		if err != nil {
			return fmt.Errorf("creating watermarks: %w", err)
		}
		watermarks[i+1] = w
	}
	err = api.AddWatermarksMapFile(path, "", watermarks, nil)
	if err != nil {
		return fmt.Errorf("adding pagenumbers: %w", err)
	}

	return nil
}

func AddStampToPdf(tp temp.TempHelper, sourcePDF string, targetPDF string, text string, offset string) {
	annotationFileName := tp.GetCompleteFileName("pdf-annotation.png")
	createWaterMarkFile(annotationFileName, text, nil)
	desc := fmt.Sprintf("scale:0.75 rel, %s, rot:0", offset)
	err := api.AddImageWatermarksFile(tp.GetCompleteFileName(sourcePDF), tp.GetCompleteFileName(targetPDF), []string{"l"}, false, annotationFileName, desc, nil)
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ImageOperationError))
	}
}

func createWaterMarkFile(fileName string, text string, customFontFile *string) {
	img := image.NewRGBA(image.Rect(0, 0, 800, 800))
	addLabel(img, 60, 60, text, customFontFile)
	f, err := os.Create(fileName)
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ImageOperationError))
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ImageOperationError))
	}
}

func addLabel(img *image.RGBA, x, y int, label string, customFontFile *string) {
	col := color.RGBA{40, 40, 140, 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	fontFile := conf.Config.Server.GetDocumentResourceFilePath(conf.Config.Fonts.Stamp)
	if customFontFile != nil {
		fontFile = *customFontFile
	}
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ResourceMissing))
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ResourceMissing))
		return
	}

	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(col),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    12,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
		Dot: point,
	}

	rows := strings.Split(label, "\n")
	lineHeight := d.Face.Metrics().Height.Ceil()
	imgHeight := img.Bounds().Max.Y
	for _, row := range rows {
		if y+lineHeight > imgHeight {
			fmt.Println("Out of bounds")
			break // Stop drawing if next line is out of bounds
		}
		d.Dot.Y = fixed.I(y) // Setze die Y-Position zurück auf den Anfang
		d.Dot.X = point.X
		d.DrawString(row)
		y += lineHeight
	}
}
