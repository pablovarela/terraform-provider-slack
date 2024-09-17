package slack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "email"},
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "email"},
			},
			"include_deactivated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to `true` to search deactivated users.",
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*slack.Client)

	var user *slack.User
	if name, ok := d.GetOk("name"); ok {
		u, err := searchByName(ctx, name.(string), client)
		if err != nil {
			return diag.FromErr(fmt.Errorf("not found %s: %w", name.(string), err))
		}
		user = u
	}

	if email, ok := d.GetOk("email"); ok {
		u, err := client.GetUserByEmailContext(ctx, email.(string))

		if err != nil {
			// We use searchByEmail (slower) to find deactivated users in special cases
			// and when asked for with `include_deactivated`
			if d.Get("include_deactivated").(bool) && err.Error() == "users_not_found" {
				u, err := searchByEmail(ctx, email.(string), client)
				if err != nil {
					return diag.FromErr(fmt.Errorf("not found %s: %w", email.(string), err))
				}
				user = u
			} else {
				return diag.FromErr(fmt.Errorf("not found %s: %w", email.(string), err))
			}
		} else {
			user = u
		}
	}

	if user == nil {
		return diag.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	d.SetId(user.ID)
	if err := d.Set("name", user.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}

	if err := d.Set("email", user.Profile.Email); err != nil {
		return diag.Errorf("error setting name: %s", err)
	}

	return diags
}

func searchByName(ctx context.Context, name string, client *slack.Client) (*slack.User, error) {
	users, err := client.GetUsersContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get workspace users: %s", err)
	}

	var matchingUsers []slack.User
	for _, user := range users {
		if user.Name == name {
			matchingUsers = append(matchingUsers, user)
		}
	}

	if len(matchingUsers) < 1 {
		return nil, fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(matchingUsers) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return &matchingUsers[0], nil
}

// searchByEmail requires listing all users and is very slow.
// This should only be used a failback for finding deactivated users.
func searchByEmail(ctx context.Context, email string, client *slack.Client) (*slack.User, error) {
	users, err := client.GetUsersContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get workspace users: %s", err)
	}

	var matchingUsers []slack.User
	for _, user := range users {
		if user.Profile.Email == email {
			matchingUsers = append(matchingUsers, user)
		}
	}

	if len(matchingUsers) < 1 {
		return nil, fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(matchingUsers) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return &matchingUsers[0], nil
}
