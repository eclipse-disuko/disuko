// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package pdf

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/logy"
)

var singletonMutex sync.Mutex

func ConvertAndMerge(rs *logy.RequestSession, outfile string, pageHeader, pageFooter *string, contentFiles []string) error {
	singletonMutex.Lock()
	defer singletonMutex.Unlock()

	var uniteArgs []string
	for _, file := range contentFiles {
		pdfPath := fmt.Sprintf("%s.pdf", file)

		headerTemplate := "<span></span>"
		if pageHeader != nil {
			headerTemplate = *pageHeader
		}
		footerTemplate := "<span></span>"
		if pageFooter != nil {
			footerTemplate = *pageFooter
		}

		uniteArgs = append(uniteArgs, pdfPath)
		ctx, cancelChrome := setupChrome()
		err := func() error {
			defer cancelChrome()

			setUserAgent := chromedp.Func(func(ctx context.Context, t *chromedp.Target) error {
				_, err := cdp.Call(ctx, t, emulation.SetUserAgentOverride, emulation.SetUserAgentOverrideParams{
					UserAgent: "Disco Browser",
				})
				return err
			})

			if err := chromedp.Do(ctx,
				setUserAgent,
				chromedp.Navigate("file://"+file),
				chromedp.WaitVisible(`body`),
			); err != nil {
				return err
			}

			buf, err := chromedp.Run(ctx, chromedp.PrintToPDF(
				chromedp.PDFPrintBackground(),
				chromedp.PDFPreferCSSPageSize(),
				chromedp.PDFHeaderTemplate(headerTemplate),
				chromedp.PDFFooterTemplate(footerTemplate),
			))
			if err != nil {
				logy.Errorf(rs, "printing %s to pdf failed: %s", file, pdfPath, err.Error())
				return err
			}
			if err := os.WriteFile(pdfPath, buf, 0644); err != nil {
				logy.Errorf(rs, "writing %s to pdf %s failed: %s", file, pdfPath, err.Error())
				return err
			}
			return nil
		}()
		if err != nil {
			return fmt.Errorf("running chromedp: %w", err)
		}
	}

	if len(uniteArgs) == 1 {
		os.Rename(uniteArgs[0], outfile)
		return nil
	}

	uniteArgs = append(uniteArgs, outfile)
	cmd := exec.Command("pdfunite", uniteArgs...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting pdfunite: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		logy.Errorw(rs, err.Error())
		return fmt.Errorf("waiting for pdfunite: %w", err)
	}
	return nil
}

func setupChrome() (context.Context, func()) {
	opts := chromedp.DefaultExecAllocatorOptions[:]
	if conf.Config.Server.VanillaDisuko {
		// no sandbox needed for chromedp when using docker compose
		opts = append(opts,
			chromedp.NoSandbox,
		)
	}
	opts = append(opts, chromedp.ProxyServer(conf.Config.Proxy.HttpProxy))
	opts = append(opts, chromedp.DisableGPU)
	ctx, cancelAllo := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
		chromedp.WithErrorf(log.Printf),
		// // TODO: delete after 2 weeks
		// chromedp.WithDebugf(log.Printf),
	)
	return ctx, func() {
		cancel()
		cancelAllo()
	}
}
