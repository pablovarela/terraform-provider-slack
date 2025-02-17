package slack

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type testUser struct {
	id    string
	name  string
	email string
}

var (
	testUserCreator = testUser{
		id:    "U01D6L97N0M",
		name:  "contact",
		email: "contact@tapih.co.uk",
	}

	testUser00 = testUser{
		id:    "U01D31S1GUE",
		name:  "contact_test-user-ter",
		email: "contact+test-user-terraform-provider-slack-00@tapih.co.uk",
	}

	testUser01 = testUser{
		id:    "U01DZK10L1W",
		name:  "contact_test-user-206",
		email: "contact+test-user-terraform-provider-slack-01@tapih.co.uk",
	}
)

var (
	testAccProvider          *schema.Provider
	testAccProviderFactories func(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error)
)

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = func(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error) {
		providerNames := []string{"slack"}
		factories := make(map[string]func() (*schema.Provider, error), len(providerNames))
		for _, name := range providerNames {
			p := testAccProvider
			factories[name] = func() (*schema.Provider, error) {
				return p, nil
			}
			*providers = append(*providers, p)
		}
		return factories
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(_ *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SLACK_TOKEN"); v == "" {
		t.Fatal("SLACK_TOKEN must be set for acceptance tests")
	}
}
