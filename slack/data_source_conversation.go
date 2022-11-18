package slack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceConversation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlackConversationRead,

		Schema: map[string]*schema.Schema{
			"channel_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purpose": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_ext_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_org_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_general": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSlackConversationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)
	channelID := d.Get("channel_id").(string)
	channelName := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	var channel *slack.Channel
	var err error
	if channelID != "" {
		channel, err = client.GetConversationInfoContext(ctx, channelID, false)
		if err != nil {
			return diag.FromErr(fmt.Errorf("couldn't get conversation info for %s: %w", channelID, err))
		}
	} else if channelName != "" {
		channel, err = findExistingChannel(ctx, client, channelName, isPrivate)
		if err != nil {
			return diag.FromErr(fmt.Errorf("couldn't get conversation info for %s: %w", channelName, err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("channel_id or name must be set"))
	}

	users, _, err := client.GetUsersInConversationContext(ctx, &slack.GetUsersInConversationParameters{
		ChannelID: channel.ID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("couldn't get users in conversation for %s: %w", channel.ID, err))
	}
	return updateChannelData(d, channel, users)
}
