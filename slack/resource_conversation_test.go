package slack

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
)

const (
	conversationNamePrefix = "test-acc-slack-conversation-test"
)

func init() {
	resource.AddTestSweepers("slack_conversation", &resource.Sweeper{
		Name: "slack_conversation",
		F: func(string) error {
			client, err := sharedSlackClient()
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}
			c := client.(*slack.Client)
			channels, _, err := c.GetConversations(&slack.GetConversationsParameters{
				ExcludeArchived: true,
				Types:           []string{"public_channel", "private_channel"},
			})
			if err != nil {
				return fmt.Errorf("[ERROR] error getting channels: %s", err)
			}
			var sweeperErrs *multierror.Error
			for _, channel := range channels {
				if strings.HasPrefix(channel.Name, conversationNamePrefix) {
					err := c.ArchiveConversationContext(context.Background(), channel.ID)
					if err != nil {
						if err.Error() != "already_archived" {
							sweeperErr := fmt.Errorf("archiving channel %s during sweep: %s", channel.Name, err)
							log.Printf("[ERROR] %s", sweeperErr)
							sweeperErrs = multierror.Append(sweeperErrs, err)
						}
					}
					fmt.Printf("[INFO] archived channel %s during sweep\n", channel.Name)
				}
			}
			return sweeperErrs.ErrorOrNil()
		},
	})
}

func TestAccSlackConversationTest(t *testing.T) {
	t.Parallel()

	resourceName := "slack_conversation.test"

	t.Run("update name, topic and purpose", func(t *testing.T) {
		name := acctest.RandomWithPrefix(conversationNamePrefix)
		createChannel := testAccSlackConversation(name)

		updateName := acctest.RandomWithPrefix(conversationNamePrefix)
		updateChannel := testAccSlackConversation(updateName)

		testSlackConversationUpdate(t, resourceName, createChannel, &updateChannel)
	})

	t.Run("archive channel", func(t *testing.T) {
		name := acctest.RandomWithPrefix(conversationNamePrefix)
		createChannel := testAccSlackConversationWithMembers(name, []string{testUser00.id})

		updateChannel := createChannel
		updateChannel.IsArchived = true

		testSlackConversationUpdate(t, resourceName, createChannel, &updateChannel)
	})

	t.Run("unarchive channel", func(t *testing.T) {
		name := acctest.RandomWithPrefix(conversationNamePrefix)
		createChannel := testAccSlackConversationWithMembers(name, []string{testUser00.id})
		createChannel.IsArchived = true

		updateChannel := createChannel
		updateChannel.IsArchived = false

		testSlackConversationUpdate(t, resourceName, createChannel, &updateChannel)
	})

	t.Run("add permanent members", func(t *testing.T) {
		name := acctest.RandomWithPrefix(conversationNamePrefix)
		createChannel := testAccSlackConversationWithMembers(name, []string{testUser00.id})

		updateChannel := createChannel
		updateChannel.Members = []string{testUser00.id, testUser01.id}

		testSlackConversationUpdate(t, resourceName, createChannel, &updateChannel)
	})

	t.Run("remove permanent members", func(t *testing.T) {
		name := acctest.RandomWithPrefix(conversationNamePrefix)
		createChannel := testAccSlackConversationWithMembers(name, []string{testUser00.id, testUser01.id})

		updateChannel := createChannel
		updateChannel.Members = []string{testUser00.id}

		testSlackConversationUpdate(t, resourceName, createChannel, &updateChannel)
	})

	t.Run("invite only the creator to the channel", func(t *testing.T) {
		name := acctest.RandomWithPrefix(conversationNamePrefix)
		users := []string{testUserCreator.id}
		createChannel := testAccSlackConversationWithMembers(name, users)

		testSlackConversationUpdate(t, resourceName, createChannel, nil)
	})

	t.Run("invite creator and other users to the channel", func(t *testing.T) {
		name := acctest.RandomWithPrefix(conversationNamePrefix)
		users := []string{testUserCreator.id, testUser00.id, testUser01.id}
		createChannel := testAccSlackConversationWithMembers(name, users)

		testSlackConversationUpdate(t, resourceName, createChannel, nil)
	})
}

func testSlackConversationUpdate(t *testing.T, resourceName string, createChannel slack.Channel, updateChannel *slack.Channel) {
	var providers []*schema.Provider
	steps := []resource.TestStep{
		{
			Config: testAccSlackConversationConfig(createChannel),
			Check: resource.ComposeTestCheckFunc(
				testCheckSlackChannelAttributes(t, resourceName, createChannel),
				testCheckResourceAttrBasic(resourceName, createChannel),
			),
		},
		{
			ResourceName:            resourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"permanent_members", "action_on_destroy", "action_on_update_permanent_members", "adopt_existing_channel"},
		},
	}

	if updateChannel != nil {
		steps = append(steps, resource.TestStep{
			Config: testAccSlackConversationConfig(*updateChannel),
			Check: resource.ComposeTestCheckFunc(
				testCheckSlackChannelAttributes(t, resourceName, *updateChannel),
				testCheckResourceAttrBasic(resourceName, *updateChannel),
			),
		},
		)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      testAccCheckConversationDestroy,
		Steps:             steps,
	})
}

func testCheckResourceAttrBasic(resourceName string, channel slack.Channel) resource.TestCheckFunc {
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
	)
}

func testCheckSlackChannelAttributes(t *testing.T, resourceName string, expectedChannel slack.Channel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		c := testAccProvider.Meta().(*slack.Client)
		primary := rs.Primary
		channel, err := c.GetConversationInfo(primary.ID, false)
		if err != nil {
			return fmt.Errorf("couldn't get conversation info for %s: %s", primary.ID, err)
		}

		require.Equal(t, primary.Attributes["name"], channel.Name, "channel name does not match")
		require.Equal(t, primary.Attributes["topic"], channel.Topic.Value, "channel topic does not match")
		require.Equal(t, primary.Attributes["purpose"], channel.Purpose.Value, "channel purpose does not match")
		require.Equal(t, primary.Attributes["creator"], channel.Creator, "channel creator does not match")
		require.Equal(t, primary.Attributes["is_private"], fmt.Sprintf("%t", channel.IsPrivate), "channel is_private does not match")
		require.Equal(t, primary.Attributes["is_archived"], fmt.Sprintf("%t", channel.IsArchived), "channel is_archived does not match")
		require.Equal(t, primary.Attributes["is_shared"], fmt.Sprintf("%t", channel.IsShared), "channel is_shared does not match")
		require.Equal(t, primary.Attributes["is_org_shared"], fmt.Sprintf("%t", channel.IsOrgShared), "channel is_org_shared does not match")
		require.Equal(t, primary.Attributes["is_ext_shared"], fmt.Sprintf("%t", channel.IsExtShared), "channel is_ext_shared does not match")
		require.Equal(t, primary.Attributes["is_general"], fmt.Sprintf("%t", channel.IsGeneral), "channel is_general does not match")

		channelUsers, _, err := c.GetUsersInConversationContext(context.Background(), &slack.GetUsersInConversationParameters{
			ChannelID: channel.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't get users in conversation for %s: %s", channel.ID, err)
		}
		definedMembers := expectedChannel.Members
		assertUsersInStateAreInTheChannel(t, primary, definedMembers, channelUsers)

		return nil
	}
}

func assertUsersInStateAreInTheChannel(t *testing.T, primary *terraform.InstanceState, definedMembers []string, users []string) {
	permanentUsersLength, _ := strconv.Atoi(primary.Attributes["permanent_members.#"])
	require.Equal(t, len(definedMembers), permanentUsersLength, "defined members length should match state")

	definedMembersPlusCreator := definedMembers
	if !contains(definedMembers, testUserCreator.id) {
		definedMembersPlusCreator = append(definedMembers, testUserCreator.id)
	}
	require.Equal(t, len(definedMembersPlusCreator), len(users), "defined members length should match users in channel")

	for i := 0; i < permanentUsersLength; i++ {
		user := primary.Attributes[fmt.Sprintf("permanent_members.%d", i)]
		require.True(t, contains(users, user), "user should be in the channel")
		require.True(t, contains(definedMembers, user), "member in state should be defined in the resource")
	}
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
