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
	namePrefix := "test-acc-slack-conversation-test"

	t.Run("test update name, topic and purpose", func(t *testing.T) {
		name := acctest.RandomWithPrefix(namePrefix)
		createChannel := testAccSlackConversation(name)

		updateName := acctest.RandomWithPrefix(fmt.Sprintf("%s-update", namePrefix))
		updateChannel := testAccSlackConversation(updateName)
		updateChannel.ID = createChannel.ID

		testSlackConversationUpdate(t, resourceName, createChannel, updateChannel)
	})

	t.Run("test archive channel", func(t *testing.T) {
		name := acctest.RandomWithPrefix(namePrefix)
		createChannel := testAccSlackConversationWithMembers(name, []string{testUser00.id})

		updateChannel := createChannel
		updateChannel.IsArchived = true

		testSlackConversationUpdate(t, resourceName, createChannel, updateChannel)
	})

	t.Run("test unarchive channel", func(t *testing.T) {
		name := acctest.RandomWithPrefix(namePrefix)
		createChannel := testAccSlackConversationWithMembers(name, []string{testUser00.id})
		createChannel.IsArchived = true

		updateChannel := createChannel
		updateChannel.IsArchived = false

		testSlackConversationUpdate(t, resourceName, createChannel, updateChannel)
	})

	//t.Run("test update permanent members", func(t *testing.T) {
	//	name := acctest.RandomWithPrefix(namePrefix)
	//	createChannel := testAccSlackConversationWithMembers(name, []string{testUser00.id})
	//
	//	updateChannel := createChannel
	//	updateChannel.Members = []string{testUser00.id, testUser01.id}
	//
	//	testSlackConversationUpdate(t, resourceName, createChannel, updateChannel)
	//})
}

func testSlackConversationUpdate(t *testing.T, resourceName string, createChannel slack.Channel, updateChannel slack.Channel) {
	var providers []*schema.Provider
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      testAccCheckConversationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlackConversationConfig(createChannel),
				Check:  testCheckResourceAttrBasic(resourceName, createChannel),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permanent_members"},
			},
			{
				Config: testAccSlackConversationConfig(updateChannel),
				Check:  testCheckResourceAttrBasic(resourceName, updateChannel),
			},
		},
	})
}

func testCheckResourceAttrBasic(resourceName string, channel slack.Channel) resource.TestCheckFunc {
	members := append(channel.Members, testUserCreator.id)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", channel.Name),
		resource.TestCheckResourceAttr(resourceName, "topic", channel.Topic.Value),
		resource.TestCheckResourceAttr(resourceName, "purpose", channel.Purpose.Value),
		resource.TestCheckResourceAttr(resourceName, "creator", testUserCreator.id),
		resource.TestCheckResourceAttr(resourceName, "is_private", fmt.Sprintf("%t", channel.IsPrivate)),
		resource.TestCheckResourceAttr(resourceName, "is_archived", fmt.Sprintf("%t", channel.IsArchived)),
		resource.TestCheckResourceAttr(resourceName, "is_shared", fmt.Sprintf("%t", channel.IsShared)),
		resource.TestCheckResourceAttr(resourceName, "is_org_shared", fmt.Sprintf("%t", channel.IsOrgShared)),
		resource.TestCheckResourceAttr(resourceName, "is_ext_shared", fmt.Sprintf("%t", channel.IsExtShared)),
		resource.TestCheckResourceAttr(resourceName, "is_general", fmt.Sprintf("%t", channel.IsGeneral)),
		testCheckResourceAttrSlice(resourceName, "permanent_members", channel.Members),
		testCheckResourceAttrSlice(resourceName, "members", members),
	)
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

func testAccSlackConversation(channelName string) slack.Channel {
	return testAccSlackConversationWithMembers(channelName, []string{})
}

func testAccSlackConversationWithMembers(channelName string, members []string) slack.Channel {
	sort.Strings(members)
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
  is_archived       = %t
}
`, c.Name, c.Topic.Value, c.Purpose.Value, strings.Join(members, ","), c.IsPrivate, c.IsArchived)
}
