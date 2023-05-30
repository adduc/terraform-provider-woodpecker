package internal

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func NewProto6ProviderFactory() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"woodpecker": providerserver.NewProtocol6WithError(New()),
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("WOODPECKER_SERVER"); v == "" {
		t.Fatal("WOODPECKER_SERVER must be set for acceptance tests")
	}

	if v := os.Getenv("WOODPECKER_TOKEN"); v == "" {
		t.Fatal("WOODPECKER_TOKEN must be set for acceptance tests")
	}
}
