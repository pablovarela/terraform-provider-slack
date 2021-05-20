package slack

import (
	"context"
	"fmt"

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
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var group *slack.UserGroup

	if name, ok := d.GetOk("name"); ok {
		u, err := findUserGroupByName(ctx, name.(string), false, m)
		if err != nil {
			return diag.FromErr(err)
		}
		group = &u
	}

	if id, ok := d.GetOk("usergroup_id"); ok {
		u, err := findUserGroupByID(ctx, id.(string), false, m)
		if err != nil {
			return diag.FromErr(err)
		}
		group = &u
	}

	if group == nil {
		return diag.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	d.SetId(group.ID)
	if err := d.Set("usergroup_id", group.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting usergroup ID: %s", err))
	}
	return updateUserGroupData(d, *group)
}
