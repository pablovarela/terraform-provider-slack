package slack

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
	"log"
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
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
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
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func configureSlackConversation(d *schema.ResourceData, channel *slack.Channel) {
	d.SetId(channel.ID)
	_ = d.Set("name", channel.Name)
	_ = d.Set("topic", channel.Topic.Value)
	_ = d.Set("purpose", channel.Purpose.Value)
	_ = d.Set("is_archived", channel.IsArchived)
	_ = d.Set("is_shared", channel.IsShared)
	_ = d.Set("is_ext_shared", channel.IsExtShared)
	_ = d.Set("is_org_shared", channel.IsOrgShared)
	_ = d.Set("created", channel.Created)
	_ = d.Set("creator", channel.Creator)

	// Required
	_ = d.Set("is_private", channel.IsPrivate)

	// Never support
	//_ = d.Set("members", channel.Members)
	//_ = d.Set("num_members", channel.NumMembers)
	//_ = d.Set("unread_count", channel.UnreadCount)
	//_ = d.Set("unread_count_display", channel.UnreadCountDisplay)
	//_ = d.Set("last_read", channel.Name)
	//_ = d.Set("latest", channel.Name)
}

func resourceSlackConversationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	log.Printf("[DEBUG] Creating Conversation: %s", name)
	channel, err := client.CreateConversationContext(ctx, name, isPrivate)

	if err != nil {
		return diag.FromErr(err)
	}

	iMembers := d.Get("members").(*schema.Set)
	if iMembers.Len() != 0 {
		userIds := make([]string, len(iMembers.List()))
		for i, v := range iMembers.List() {
			userIds[i] = v.(string)
		}
		_, err = client.InviteUsersToConversation(channel.ID, userIds...)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if topic, ok := d.GetOk("topic"); ok {
		if _, err := client.SetTopicOfConversationContext(ctx, channel.ID, topic.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	configureSlackConversation(d, channel)

	return resourceSlackConversationRead(ctx, d, m)
}

func resourceSlackConversationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*slack.Client)

	id := d.Id()

	log.Printf("[DEBUG] Reading Conversation: %s", d.Id())
	channel, err := client.GetConversationInfoContext(ctx, id, false)

	if err != nil {
		return diag.FromErr(err)
	}

	configureSlackConversation(d, channel)

	return diags
}

func resourceSlackConversationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)

	id := d.Id()

	if _, err := client.RenameConversationContext(ctx, id, d.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	}

	if topic, ok := d.GetOk("topic"); ok {
		if _, err := client.SetTopicOfConversationContext(ctx, id, topic.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if purpose, ok := d.GetOk("purpose"); ok {
		if _, err := client.SetPurposeOfConversationContext(ctx, id, purpose.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if isArchived, ok := d.GetOk("is_archived"); ok {
		if isArchived.(bool) {
			if err := client.ArchiveConversationContext(ctx, id); err != nil {
				if err.Error() != "already_archived" {
					return diag.FromErr(err)
				}
			}
		} else {
			if err := client.UnArchiveConversationContext(ctx, id); err != nil {
				if err.Error() != "not_archived" {
					return diag.FromErr(err)
				}
			}
		}
	}

	return resourceSlackConversationRead(ctx, d, m)
}

func resourceSlackConversationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*slack.Client)

	id := d.Id()
	log.Printf("[DEBUG] Deleting(archive) Conversation: %s (%s)", id, d.Get("name"))
	if err := client.ArchiveConversationContext(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
