package slack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupsRead,

		Schema: map[string]*schema.Schema{
			"usergroups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceUserGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SlackClient)

	userGroups, err := client.GetUserGroupsCached(ctx, false, true)
	if err != nil {
		return diag.Errorf("couldn't get usergroups: %s", err)
	}

	groups := make([]map[string]interface{}, len(userGroups))
	for i, ug := range userGroups {
		groups[i] = map[string]interface{}{
			"id":          ug.ID,
			"name":        ug.Name,
			"handle":      ug.Handle,
			"description": ug.Description,
			"channels":    ug.Prefs.Channels,
			"users":       ug.Users,
		}
	}

	d.SetId("usergroups")
	if err := d.Set("usergroups", groups); err != nil {
		return diag.Errorf("error setting usergroups: %s", err)
	}

	return nil
}
