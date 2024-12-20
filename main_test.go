package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"strings"
	"testing"
	"time"
)

// Helper function to create a ChromeDP context
func createChromeDPContext(t *testing.T, headless, disableGPU, startMaximized bool) (context.Context, context.CancelFunc) {
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts,
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", disableGPU),
		chromedp.Flag("start-maximized", startMaximized),
		chromedp.Flag("no-sandbox", true),
	)
	allocCtx, cancelAllocator := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)

	// Handle cleanup
	cancelAll := func() {
		cancelTimeout()
		cancelCtx()
		cancelAllocator()
	}

	return timeoutCtx, cancelAll
}

// Test that the search results contain "Golang"
func TestGoogleSearch_Golang(t *testing.T) {
	ctx, cancel := createChromeDPContext(t, true, true, false)
	defer cancel()

	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.google.com"),
		chromedp.WaitVisible(`//textarea[@name="q"]`),
		chromedp.SendKeys(`//textarea[@name="q"]`, "Golang"),
		chromedp.SendKeys(`//textarea[@name="q"]`, kb.Enter),
		chromedp.WaitVisible(`#search`),
		chromedp.Text(`#search`, &result),
	)

	if err != nil {
		t.Fatalf("Test Failed: %v", err)
	}

	if strings.Contains(result, "Golang") {
		t.Log("Test Passed: 'Golang' found in search results")
	} else {
		t.Errorf("Test Failed: 'Golang' not found in search results")
	}
}

//// Test that the search bar is visible
//func TestGoogleSearch_SearchBarVisible(t *testing.T) {
//	ctx, cancel := createChromeDPContext(t, false, false, true)
//	defer cancel()
//
//	err := chromedp.Run(ctx,
//		chromedp.Navigate("https://www.google.com"),
//		chromedp.WaitVisible(`//textarea[@name="q"]`),
//	)
//
//	if err != nil {
//		t.Errorf("Test Failed: Search bar not visible - %v", err)
//	} else {
//		t.Log("Test Passed: Search bar is visible")
//	}
//}
//
//// Test that there are no results when the search query is empty
//func TestGoogleSearch_EmptyQuery(t *testing.T) {
//	ctx, cancel := createChromeDPContext(t, false, false, true)
//	defer cancel()
//
//	var result string
//	err := chromedp.Run(ctx,
//		chromedp.Navigate("https://www.google.com"),
//		chromedp.WaitVisible(`//textarea[@name="q"]`),
//		chromedp.SendKeys(`//textarea[@name="q"]`, ""), // Empty query
//		chromedp.SendKeys(`//textarea[@name="q"]`, kb.Enter),
//		chromedp.WaitVisible(`#search`),
//		chromedp.Text(`#search`, &result),
//	)
//
//	if err != nil {
//		t.Fatalf("Test Failed: %v", err)
//	}
//
//	if result == "" {
//		t.Log("Test Passed: No search results for empty query")
//	} else {
//		t.Errorf("Test Failed: Unexpected results for empty query")
//	}
//}
