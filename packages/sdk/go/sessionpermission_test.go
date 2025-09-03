// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kenkaiii/ken8n-coder-sdk-go"
	"github.com/kenkaiii/ken8n-coder-sdk-go/internal/testutil"
	"github.com/kenkaiii/ken8n-coder-sdk-go/option"
)

func TestSessionPermissionRespond(t *testing.T) {
	t.Skip("skipped: tests are disabled for the time being")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := opencode.NewClient(
		option.WithBaseURL(baseURL),
	)
	_, err := client.Session.Permissions.Respond(
		context.TODO(),
		"id",
		"permissionID",
		opencode.SessionPermissionRespondParams{
			Response: opencode.F(opencode.SessionPermissionRespondParamsResponseOnce),
		},
	)
	if err != nil {
		var apierr *opencode.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
