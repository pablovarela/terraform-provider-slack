package slack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/slack-go/slack"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestAccSlackConversationTest(t *testing.T) {
	t.Parallel()

	resourceName := "slack_conversation.test"

	name := acctest.RandomWithPrefix("test-acc-slack-conversation-test")
	var permanentMembers []string
	var expectedMembers = []string{testUserCreator.id}
	createChannel := testAccSlackConversation(name, permanentMembers)

	var updatedPermanentMembers = []string{testUser00.id}
	sort.Strings(updatedPermanentMembers)
	var updatedExpectedMembers = []string{testUserCreator.id, testUser00.id}
	sort.Strings(updatedExpectedMembers)
	updateName := acctest.RandomWithPrefix("test-acc-slack-conversation-test-update")
	updateChannel := testAccSlackConversation(updateName, updatedPermanentMembers)
	updateChannel.ID = createChannel.ID

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
				Config: testAccSlackConversationConfig(createChannel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", createChannel.Name),
					resource.TestCheckResourceAttr(resourceName, "topic", createChannel.Topic.Value),
					resource.TestCheckResourceAttr(resourceName, "purpose", createChannel.Purpose.Value),
					resource.TestCheckResourceAttr(resourceName, "creator", testUserCreator.id),
					resource.TestCheckResourceAttr(resourceName, "is_private", fmt.Sprintf("%t", createChannel.IsPrivate)),
					resource.TestCheckResourceAttr(resourceName, "is_archived", fmt.Sprintf("%t", createChannel.IsArchived)),
					resource.TestCheckResourceAttr(resourceName, "is_shared", fmt.Sprintf("%t", createChannel.IsShared)),
					resource.TestCheckResourceAttr(resourceName, "is_org_shared", fmt.Sprintf("%t", createChannel.IsOrgShared)),
					resource.TestCheckResourceAttr(resourceName, "is_ext_shared", fmt.Sprintf("%t", createChannel.IsExtShared)),
					resource.TestCheckResourceAttr(resourceName, "is_general", fmt.Sprintf("%t", createChannel.IsGeneral)),
					testCheckResourceAttrSlice(resourceName, "permanent_members", permanentMembers),
					testCheckResourceAttrSlice(resourceName, "members", expectedMembers),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permanent_members"},
			},
			{
				Config: testAccSlackConversationConfig(updateChannel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateChannel.Name),
					resource.TestCheckResourceAttr(resourceName, "topic", updateChannel.Topic.Value),
					resource.TestCheckResourceAttr(resourceName, "purpose", updateChannel.Purpose.Value),
					resource.TestCheckResourceAttr(resourceName, "creator", testUserCreator.id),
					resource.TestCheckResourceAttr(resourceName, "is_private", fmt.Sprintf("%t", updateChannel.IsPrivate)),
					resource.TestCheckResourceAttr(resourceName, "is_archived", fmt.Sprintf("%t", updateChannel.IsArchived)),
					resource.TestCheckResourceAttr(resourceName, "is_shared", fmt.Sprintf("%t", updateChannel.IsShared)),
					resource.TestCheckResourceAttr(resourceName, "is_org_shared", fmt.Sprintf("%t", updateChannel.IsOrgShared)),
					resource.TestCheckResourceAttr(resourceName, "is_ext_shared", fmt.Sprintf("%t", updateChannel.IsExtShared)),
					resource.TestCheckResourceAttr(resourceName, "is_general", fmt.Sprintf("%t", updateChannel.IsGeneral)),
					testCheckResourceAttrSlice(resourceName, "permanent_members", updatedPermanentMembers),
					//testCheckResourceAttrSlice(resourceName, "members", updatedExpectedMembers),
				),
			},
		},
	})
}

func testCheckResourceAttrSlice(resourceName string, key string, a []string) resource.TestCheckFunc {
	tests := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.#", key), strconv.Itoa(len(a))),
	}

	for i, v := range a {
		tests = append(
			tests,
			resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("%s.%d", key, i), v),
		)
	}

	return resource.ComposeTestCheckFunc(tests...)
}

func testAccSlackConversation(channelName string, members []string) slack.Channel {
	channel := slack.Channel{
		GroupConversation: slack.GroupConversation{
			Name: channelName,
			Topic: slack.Topic{
				Value: fmt.Sprintf("Topic for %s", channelName),
			},
			Purpose: slack.Purpose{
				Value: fmt.Sprintf("Purpose of %s", channelName),
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
		if rs.Type != "slack_conversation" {
			continue
		}

		err := archiveConversationWithContext(context.Background(), c, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error archiving channel %s: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccSlackConversationConfig(c slack.Channel) string {
	var members []string
	for _, member := range c.Members {
		members = append(members, fmt.Sprintf(`"%s"`, member))
	}

	return fmt.Sprintf(`
resource slack_conversation test {
  name              = "%s"
  topic             = "%s"
  purpose           = "%s"
  permanent_members = [%s]
  is_private        = %t
}
`, c.Name, c.Topic.Value, c.Purpose.Value, strings.Join(members, ","), c.IsPrivate)
}
