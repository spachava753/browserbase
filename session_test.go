package browserbase

import (
	"context"
	"os"
	"testing"
)

// Implements SecuritySource for ogen.
type envAPIKeySource struct {
	key string
}

func (s envAPIKeySource) XBBAPIKey(ctx context.Context, operationName OperationName) (XBBAPIKey, error) {
	return XBBAPIKey{APIKey: s.key}, nil
}

func testClientAndProjectID(t *testing.T) (*Client, string) {
	t.Helper()
	apiKey := os.Getenv("X_BB_API_Key")
	if apiKey == "" {
		t.Fatal("X_BB_API_Key environment variable not set")
	}
	projectID := os.Getenv("BB_PROJECT_ID")
	if projectID == "" {
		t.Fatal("BB_PROJECT_ID environment variable not set")
	}

	sec := &EnvAPIKeySource{APIKey: apiKey}
	client, err := NewClient("https://api.browserbase.com", sec)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client, projectID
}

func TestCreateSession(t *testing.T) {
	client, projectID := testClientAndProjectID(t)
	ctx := context.Background()
	params := &SessionCreateParams{ProjectId: projectID}

	resp, err := client.CreateSession(ctx, params)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	if resp.ID == "" {
		t.Errorf("expected session ID, got empty response: %+v", resp)
	}
	t.Logf("Created session with ID: %s", resp.ID)
	t.Logf("Session connect url: %s", resp.ConnectUrl)
	t.Logf("Session web driver url: %s", resp.SeleniumRemoteUrl)
}

func TestListSessions(t *testing.T) {
	client, _ := testClientAndProjectID(t)
	ctx := context.Background()
	params := ListSessionsParams{} // No filters

	sessions, err := client.ListSessions(ctx, params)
	if err != nil {
		t.Fatalf("failed to list sessions: %v", err)
	}
	t.Logf("Found %d sessions", len(sessions))
	if len(sessions) > 0 {
		t.Logf("First session ID: %s", sessions[0].ID)
	}
}

func TestDeleteSession(t *testing.T) {
	client, projectID := testClientAndProjectID(t)
	ctx := context.Background()

	// First, create a session to ensure we have one to delete
	createResp, err := client.CreateSession(ctx, &SessionCreateParams{ProjectId: projectID})
	if err != nil {
		t.Fatalf("failed to create session for delete: %v", err)
	}
	if createResp.ID == "" {
		t.Fatal("No session ID returned from create session")
	}

	params := UpdateSessionParams{ID: createResp.ID}
	delParams := &SessionUpdateParams{ProjectId: projectID, Status: SessionUpdateParamsStatusREQUESTRELEASE}

	_, err = client.UpdateSession(ctx, delParams, params)
	if err != nil {
		t.Fatalf("failed to delete (request release) session: %v", err)
	}
	t.Logf("Requested deletion (release) of session ID: %s", createResp.ID)
}
