package browserbase

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"net/url"
	"os"
)

func Example_google() {
	apiKey := os.Getenv("X_BB_API_Key")
	if apiKey == "" {
		fmt.Println("X_BB_API_Key environment variable not set")
		return
	}
	projectID := os.Getenv("BB_PROJECT_ID")
	if projectID == "" {
		fmt.Println("BB_PROJECT_ID environment variable not set")
		return
	}
	sec := &EnvAPIKeySource{APIKey: apiKey}
	client, err := NewClient("https://api.browserbase.com", sec)
	if err != nil {
		fmt.Printf("failed to create client: %s\n", err)
		return
	}

	ctx := context.Background()
	params := &SessionCreateParams{ProjectId: projectID}

	session, err := client.CreateSession(ctx, params)
	if err != nil {
		fmt.Printf("failed to create session: %s\n", err)
		return
	}
	connectURL, err := url.Parse(session.ConnectUrl)
	if err != nil {
		fmt.Printf("error parsing connect url: %s\n", err)
		return
	}
	if connectURL == nil {
		fmt.Printf("connectUrl is empty in response: %+v", session)
		return
	}
	defer func() {
		// Clean up session (delete with REQUEST_RELEASE)
		params := UpdateSessionParams{ID: session.ID}
		delParams := &SessionUpdateParams{ProjectId: projectID, Status: SessionUpdateParamsStatusREQUESTRELEASE}
		_, err := client.UpdateSession(ctx, delParams, params)
		if err != nil {
			fmt.Printf("failed to delete (request release) session: %v", err)
		}
	}()

	// chromedp requires that the port be specified in the URL.
	connectURL = &url.URL{
		Scheme:   connectURL.Scheme,
		Host:     connectURL.Host + ":443",
		RawQuery: connectURL.RawQuery,
	}

	// Connect to the running Chrome instance in BrowserBase using chromedp Remote Allocator.
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(ctx, connectURL.String(), chromedp.NoModifyURL)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	var pageTitle string
	err = chromedp.Run(chromeCtx,
		chromedp.Navigate("https://www.google.com/"),
		// Wait for the page to load by waiting for the input box to become visible
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Title(&pageTitle),
	)
	if err != nil {
		fmt.Printf("chromedp navigation failed: %v", err)
		return
	}
	fmt.Printf("Navigated to google.com. Page title: %q", pageTitle)

	// Output: Navigated to google.com. Page title: "Google"
}
