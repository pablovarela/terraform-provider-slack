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
			"is_private": {
				Type:     schema.TypeBool,
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

	var channelID string
	if value, ok := d.GetOk("channel_id"); ok {
		channelID = value.(string)
	}

	var channelName string
	if value, ok := d.GetOk("name"); ok {
		channelName = value.(string)
	}

	if (channelID == "" && channelName == "") || (channelID != "" && channelName != "") {
		return diag.Errorf("exactly one of channel id or channel name may be specified")
	}

	var channel *slack.Channel
	var err error
	if channelID != "" {
		channel, err = client.GetConversationInfoContext(ctx, channelID, false)
		if err != nil {
			return diag.FromErr(fmt.Errorf("couldn't get conversation info for %s: %w", channelID, err))
		}
	} else {
		channel, err = findExistingChannel(ctx, client, channelName, false)
		if err != nil {
			return diag.FromErr(fmt.Errorf("couldn't get conversation info for %s: %w", channelName, err))
		}
	}

	users, _, err := client.GetUsersInConversationContext(ctx, &slack.GetUsersInConversationParameters{
		ChannelID: channel.ID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("couldn't get users in conversation for %s: %w", channel.ID, err))
	}
	return updateChannelData(d, channel, users)
}
