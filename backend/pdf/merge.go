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

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
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
		tasks := chromedp.Tasks{
			emulation.SetUserAgentOverride("Disco Browser"),
			chromedp.Navigate("file://" + file),
			chromedp.WaitVisible(`body`, chromedp.ByQuery),
			chromedp.ActionFunc(func(ctx context.Context) error {
				print := page.PrintToPDF().
					WithPrintBackground(true).
					WithPreferCSSPageSize(true).
					WithHeaderTemplate("<span></span>").
					WithFooterTemplate("<span></span>").
					WithDisplayHeaderFooter(true)

				if pageHeader != nil {
					print = print.WithHeaderTemplate(*pageHeader)
				}
				if pageFooter != nil {
					print = print.WithFooterTemplate(*pageFooter)
				}
				buf, _, err := print.Do(ctx)
				if err != nil {
					logy.Errorf(rs, "printing %s to pdf failed: %s", file, pdfPath, err.Error())
					return err
				}
				if err := os.WriteFile(pdfPath, buf, 0644); err != nil {
					logy.Errorf(rs, "writing %s to pdf %s failed: %s", file, pdfPath, err.Error())
					return err
				}
				return nil
			}),
		}
		uniteArgs = append(uniteArgs, pdfPath)
		ctx, cancelChrome := setupChrome()
		err := chromedp.Run(ctx, tasks)
		cancelChrome()
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
