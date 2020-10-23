package slack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

func TestAccSlackUserDataSource_basic(t *testing.T) {
	dataSourceName := "data.slack_user.test"

	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckSlackUserDataSourceConfigNonExistent,
				ExpectError: regexp.MustCompile(`your query returned no results`),
			},
			{
				Config: testAccCheckSlackUserDataSourceConfigExistent,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlackUserDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", testUserName),
					resource.TestCheckResourceAttr(dataSourceName, "id", testUserID),
				),
			},
		},
	})
}

func testAccCheckSlackUserDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find slack conversation data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("slack conversation data source ID not set")
		}
		return nil
	}
}

const (
	testAccCheckSlackUserDataSourceConfigNonExistent = `
data slack_user test {
 name = "non-existent"
}
`
	testAccCheckSlackUserDataSourceConfigExistent = `
data slack_user test {
 name = "pablo.varelapaz"
}
`
	testUserID                 = "U0150MARZEY"
	nonAuthenticatedTestUserID = "U014V8XMMB5"
	testUserName               = "pablo.varelapaz"
)
