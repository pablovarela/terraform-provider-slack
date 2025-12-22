package slack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupRead,

		Schema: map[string]*schema.Schema{
			"usergroup_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "usergroup_id"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "usergroup_id"},
			},
			"channels": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"handle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SlackClient)

	userGroups, err := client.GetUserGroupsCached(ctx, false, true)
	if err != nil {
		return diag.Errorf("couldn't get usergroups: %s", err)
	}

	name, hasName := d.GetOk("name")
	id, hasID := d.GetOk("usergroup_id")

	var group *slack.UserGroup
	for _, ug := range userGroups {
		if hasName && ug.Name == name.(string) {
			group = &ug
			break
		}
		if hasID && ug.ID == id.(string) {
			group = &ug
			break
		}
	}

	if group == nil {
		return diag.Errorf("usergroup not found")
	}

	d.SetId(group.ID)
	_ = d.Set("usergroup_id", group.ID)
	return updateUserGroupData(d, *group)
}
