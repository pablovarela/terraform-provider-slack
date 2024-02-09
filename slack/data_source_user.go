package slack

import (
	"context"
	"fmt"
	"time"

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
		for {
			var err error
			user, err = client.GetUserByEmailContext(ctx, email.(string))

			if rateLimitedError, ok := err.(*slack.RateLimitedError); ok {
				fmt.Printf("Rate limited. Retrying after %v seconds...\n", rateLimitedError.RetryAfter)
				time.Sleep(rateLimitedError.RetryAfter)
				continue
			} else if err != nil {
				return diag.FromErr(fmt.Errorf("not found %s: %w", email.(string), err))
			}
			break
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
