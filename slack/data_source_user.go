package slack

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*slack.Client)

	users, err := client.GetUsers()
	if err != nil {
		return diag.FromErr(err)
	}

	var matchingUsers []slack.User
	for _, user := range users {
		if user.Name == d.Get("name") {
			matchingUsers = append(matchingUsers, user)
		}
	}

	if len(matchingUsers) < 1 {
		return diag.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(matchingUsers) > 1 {
		return diag.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	user := matchingUsers[0]

	d.SetId(user.ID)
	_ = d.Set("name", user.Name)

	if err := d.Set("name", user.Name); err != nil {
		return diag.Errorf("error setting name: %s", err)
	}

	return diags
}
