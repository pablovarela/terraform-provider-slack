package slack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackConversation() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSlackConversationRead,
		CreateContext: resourceSlackConversationCreate,
		UpdateContext: resourceSlackConversationUpdate,
		DeleteContext: resourceSlackConversationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"purpose": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permanent_members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
				Optional: true,
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
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
				Required: true,
			},
			"is_archived": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func resourceSlackConversationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	channel, err := client.CreateConversationContext(ctx, name, isPrivate)
	if err != nil {
		return diag.Errorf("could create conversation %s: %s", name, err)
	}

	err = updateChannelMembers(d, client, channel.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	if topic, ok := d.GetOk("topic"); ok {
		if _, err := client.SetTopicOfConversationContext(ctx, channel.ID, topic.(string)); err != nil {
			return diag.Errorf("couldn't set conversation topic %s: %s", topic.(string), err)
		}
	}

	if purpose, ok := d.GetOk("purpose"); ok {
		if _, err := client.SetPurposeOfConversationContext(ctx, channel.ID, purpose.(string)); err != nil {
			return diag.Errorf("couldn't set conversation purpose %s: %s", purpose.(string), err)
		}
	}

	if isArchived, ok := d.GetOk("is_archived"); ok {
		if isArchived.(bool) {
			err := archiveConversationWithContext(ctx, client, channel.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(channel.ID)
	return resourceSlackConversationRead(ctx, d, m)
}

func updateChannelMembers(d *schema.ResourceData, client *slack.Client, channelID string) error {
	members := d.Get("permanent_members").(*schema.Set)
	//	fmt.Printf("[DEBUG] updating members: %d\n", members.Len())

	if members.Len() != 0 {
		userIds := schemaSetToSlice(members)
		//		fmt.Printf("[DEBUG] updating members %s\n", userIds)

		if _, err := client.InviteUsersToConversation(channelID, userIds...); err != nil {
			if err.Error() != "already_in_channel" {
				return fmt.Errorf("couldn't invite users to conversation: %s", err)
			}
		}
	}
	return nil
}

func resourceSlackConversationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)
	id := d.Id()
	return readChannelInfo(ctx, d, client, id)
}

func resourceSlackConversationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)

	id := d.Id()

	if _, err := client.RenameConversationContext(ctx, id, d.Get("name").(string)); err != nil {
		return diag.Errorf("couldn't rename conversation: %s", err)
	}

	if d.HasChange("topic") {
		topic := d.Get("topic")
		if _, err := client.SetTopicOfConversationContext(ctx, id, topic.(string)); err != nil {
			return diag.Errorf("couldn't set conversation topic %s: %s", topic.(string), err)
		}
	}

	if d.HasChange("purpose") {
		purpose := d.Get("purpose")
		if _, err := client.SetPurposeOfConversationContext(ctx, id, purpose.(string)); err != nil {
			return diag.Errorf("couldn't set conversation purpose %s: %s", purpose.(string), err)
		}
	}

	if d.HasChange("is_archived") {
		isArchived := d.Get("is_archived")
		if isArchived.(bool) {
			err := archiveConversationWithContext(ctx, client, id)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.UnArchiveConversationContext(ctx, id); err != nil {
				if err.Error() != "not_archived" {
					return diag.Errorf("couldn't archive conversation %s: %s", id, err)
				}
			}
		}
	}

	if d.HasChange("permanent_members") {
		err := updateChannelMembers(d, client, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSlackConversationRead(ctx, d, m)
}

func archiveConversationWithContext(ctx context.Context, client *slack.Client, id string) error {
	if err := client.ArchiveConversationContext(ctx, id); err != nil {
		if err.Error() != "already_archived" {
			return fmt.Errorf("couldn't archive conversation %s: %s", id, err)
		}
	}
	return nil
}

func resourceSlackConversationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*slack.Client)

	id := d.Id()
	err := archiveConversationWithContext(ctx, client, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
