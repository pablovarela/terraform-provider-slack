package slack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/slack-go/slack"
	"strings"
	"testing"
)

func TestAccSlackConversationTest(t *testing.T) {
	t.Parallel()

	resourceName := "slack_conversation.test"

	name := acctest.RandomWithPrefix("test-acc-slack-conversation-test")
	createChannel := testAccSlackConversation(name)

	updateName := acctest.RandomWithPrefix("test-acc-slack-conversation-test-update")
	updateChannel := createChannel
	updateChannel.Name = updateName

	var providers []*schema.Provider
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      testAccCheckConversationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConversationConfig(createChannel),
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckInventoryItemExists(t, resourceName, &item),
					//testAccCheckInventoryItemMatches(t, createItem, &item),
					resource.TestCheckResourceAttr(resourceName, "name", createChannel.Name),
					resource.TestCheckResourceAttr(resourceName, "topic", createChannel.Topic.Value),
					//resource.TestCheckResourceAttr(resourceName, "is_private", createChannel.IsPrivate),

				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConversationConfig(updateChannel),
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckInventoryItemExists(t, resourceName, &item),
					//testAccCheckInventoryItemMatches(t, expectedItem, &item),
					resource.TestCheckResourceAttr(resourceName, "name", updateChannel.Name),
					resource.TestCheckResourceAttr(resourceName, "topic", updateChannel.Topic.Value),
					//resource.TestCheckResourceAttr(resourceName, "is_private", createChannel.IsPrivate),
				),
			},
		},
	})
}

func testAccSlackConversation(channelName string) slack.Channel {
	var members []string
	channel := slack.Channel{
		GroupConversation: slack.GroupConversation{
			Name: channelName,
			Topic: slack.Topic{
				Value: fmt.Sprintf("Topic for %s", channelName),
			},
			Conversation: slack.Conversation{
				IsPrivate: true,
			},
			Members: members,
		},
	}
	return channel
}

func testAccCheckConversationDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*slack.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "change_inventory_item" {
			continue
		}

		err := c.ArchiveChannelContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error archiving channel %s: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccConversationConfig(c slack.Channel) string {
	var members []string
	for _, member := range c.Members {
		members = append(members, fmt.Sprintf(`"%s"`, member))
	}

	return fmt.Sprintf(`
resource slack_conversation test {
  name       = "%s"
  topic      = "%s"
  members    = [%s]
  is_private = %t
}
`, c.Name, c.Topic.Value, strings.Join(members, ","), c.IsPrivate)
}
