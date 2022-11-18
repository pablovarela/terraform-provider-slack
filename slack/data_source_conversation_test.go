package slack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/slack-go/slack"
)

func TestAccSlackConversationDataSource_basic(t *testing.T) {
	resourceName := "slack_conversation.test"
	dataSourceName := "data.slack_conversation.test"

	var providers []*schema.Provider
	name := acctest.RandomWithPrefix("test-acc-slack-conversation-test")
	members := []string{testUser00.id, testUser01.id}
	createChannel := testAccSlackConversationWithMembers(name, members)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckSlackConversationDataSourceConfigNonExistent,
				ExpectError: regexp.MustCompile(`channel_not_found`),
			},
			{
				Config: testAccCheckSlackConversationDataSourceConfig(createChannel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlackConversationDataSourceID(dataSourceName),
					resource.TestCheckResourceAttrPair(dataSourceName, "channel_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topic", resourceName, "topic"),
					resource.TestCheckResourceAttrPair(dataSourceName, "purpose", resourceName, "purpose"),
					resource.TestCheckResourceAttrPair(dataSourceName, "creator", resourceName, "creator"),
					resource.TestCheckResourceAttrPair(dataSourceName, "created", resourceName, "created"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_private", resourceName, "is_private"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_archived", resourceName, "is_archived"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_shared", resourceName, "is_shared"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_ext_shared", resourceName, "is_ext_shared"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_org_shared", resourceName, "is_org_shared"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_general", resourceName, "is_general"),
				),
			},
			{
				Config: testAccCheckSlackConversationDataSourceConfigName(createChannel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlackConversationDataSourceID(dataSourceName),
					resource.TestCheckResourceAttrPair(dataSourceName, "channel_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topic", resourceName, "topic"),
					resource.TestCheckResourceAttrPair(dataSourceName, "purpose", resourceName, "purpose"),
					resource.TestCheckResourceAttrPair(dataSourceName, "creator", resourceName, "creator"),
					resource.TestCheckResourceAttrPair(dataSourceName, "created", resourceName, "created"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_private", resourceName, "is_private"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_archived", resourceName, "is_archived"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_shared", resourceName, "is_shared"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_ext_shared", resourceName, "is_ext_shared"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_org_shared", resourceName, "is_org_shared"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_general", resourceName, "is_general"),
				),
			},
		},
	})
}

func testAccCheckSlackConversationDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find slack conversation data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("slack conversation data source id not set")
		}
		return nil
	}
}

const (
	testAccCheckSlackConversationDataSourceConfigNonExistent = `
data slack_conversation test {
 channel_id = "non-existent"
}
`
	testAccCheckSlackConversationDataSourceConfigExistent = `
data slack_conversation test {
  channel_id = slack_conversation.test.id
}
`
	testAccCheckSlackConversationDataSourceConfigNameExistent = `
data slack_conversation test {
  channel_id = slack_conversation.test.id
}
`
)

func testAccCheckSlackConversationDataSourceConfig(channel slack.Channel) string {
	return testAccSlackConversationConfig(channel) + testAccCheckSlackConversationDataSourceConfigExistent
}

func testAccCheckSlackConversationDataSourceConfigName(channel slack.Channel) string {
	return testAccSlackConversationConfig(channel) + testAccCheckSlackConversationDataSourceConfigNameExistent
}
