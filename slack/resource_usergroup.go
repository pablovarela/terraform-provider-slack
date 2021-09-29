package slack

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSlackUserGroupRead,
		CreateContext: resourceSlackUserGroupCreate,
		UpdateContext: resourceSlackUserGroupUpdate,
		DeleteContext: resourceSlackUserGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"channels": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"handle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
				Optional: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceSlackUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	handle := d.Get("handle").(string)
	channels := d.Get("channels").(*schema.Set)
	users := d.Get("users").(*schema.Set)

	userGroup := slack.UserGroup{
		Name:        name,
		Description: description,
		Handle:      handle,
		Prefs: slack.UserGroupPrefs{
			Channels: schemaSetToSlice(channels),
		},
	}
	createdUserGroup, err := client.CreateUserGroupContext(ctx, userGroup)
	if err != nil {
		if err.Error() != "name_already_exists" && err.Error() != "handle_already_exists" {
			return diag.Errorf("could not create usergroup %s: %s", name, err)
		}
		group, err := findUserGroupByName(ctx, name, true, m)
		if err != nil {
			return diag.Errorf("could not find usergroup %s: %s", name, err)
		}
		_, err = client.EnableUserGroupContext(ctx, group.ID)
		if err != nil {
			if err.Error() != "already_enabled" {
				return diag.Errorf("could not enable usergroup %s (%s): %s", name, group.ID, err)
			}
		}
		userGroup.ID = group.ID
		_, err = client.UpdateUserGroupContext(ctx, userGroup)
		if err != nil {
			return diag.Errorf("could not update usergroup %s: %s", name, err)
		}
		d.SetId(group.ID)
	} else {
		d.SetId(createdUserGroup.ID)
	}

	if users.Len() > 0 {
		_, err := client.UpdateUserGroupMembersContext(ctx, d.Id(), strings.Join(schemaSetToSlice(users), ","))
		if err != nil {
			return diag.Errorf("could not update usergroup members %s: %s", name, err)
		}
		schemaSetToSlice(users)
	}
	return resourceSlackUserGroupRead(ctx, d, m)
}

func resourceSlackUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)
	id := d.Id()

	var (
		diags      diag.Diagnostics
		userGroups []slack.UserGroup
		backoff    = &Backoff{Base: time.Second, Cap: 15 * time.Second}
	)
	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		var (
			err   error
			rlerr *slack.RateLimitedError
		)
		userGroups, err = client.GetUserGroupsContext(ctx, slack.GetUserGroupsOptionIncludeUsers(true))
		if errors.As(err, &rlerr) {
			backoff.Sleep(ctx)
			return resource.RetryableError(err)
		} else if err != nil {
			return resource.NonRetryableError(fmt.Errorf("couldn't get usergroups: %w", err))
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, userGroup := range userGroups {
		if userGroup.ID == id {
			return updateUserGroupData(d, userGroup)
		}
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("usergroup with ID %s not found, removing from state", id),
	})
	d.SetId("")
	return diags
}

type Backoff struct {
	Attempt int
	Base    time.Duration
	Cap     time.Duration
}

func (b *Backoff) Sleep(ctx context.Context) {
	b.Attempt++

	wait := b.Base * (2 << b.Attempt)
	if wait > b.Cap {
		wait = b.Cap
	}

	wait = time.Duration(rand.Int63n(int64(wait)))

	select {
	case <-time.After(wait):
	case <-ctx.Done():
	}

}

func findUserGroupByName(ctx context.Context, name string, includeDisabled bool, m interface{}) (slack.UserGroup, error) {
	client := m.(*slack.Client)
	userGroups, err := client.GetUserGroupsContext(ctx, slack.GetUserGroupsOptionIncludeDisabled(includeDisabled), slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		return slack.UserGroup{}, err
	}

	for _, userGroup := range userGroups {
		if userGroup.Name == name {
			return userGroup, nil
		}
	}

	return slack.UserGroup{}, fmt.Errorf("could not find usergroup %s", name)
}

func findUserGroupByID(ctx context.Context, id string, includeDisabled bool, m interface{}) (slack.UserGroup, error) {
	client := m.(*slack.Client)
	userGroups, err := client.GetUserGroupsContext(ctx, slack.GetUserGroupsOptionIncludeDisabled(includeDisabled), slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		return slack.UserGroup{}, err
	}

	for _, userGroup := range userGroups {
		if userGroup.ID == id {
			return userGroup, nil
		}
	}

	return slack.UserGroup{}, fmt.Errorf("could not find usergroup %s", id)
}

func resourceSlackUserGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*slack.Client)

	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	handle := d.Get("handle").(string)
	channels := d.Get("channels").(*schema.Set)
	users := d.Get("users").(*schema.Set)

	userGroup := slack.UserGroup{
		ID:          id,
		Name:        name,
		Description: description,
		Handle:      handle,
		Prefs: slack.UserGroupPrefs{
			Channels: schemaSetToSlice(channels),
		},
	}
	_, err := client.UpdateUserGroupContext(ctx, userGroup)
	if err != nil {
		return diag.Errorf("could not update usergroup %s: %s", name, err)
	}

	if d.HasChanges("users") {
		_, err := client.UpdateUserGroupMembersContext(ctx, id, strings.Join(schemaSetToSlice(users), ","))
		if err != nil {
			return diag.Errorf("could not update usergroup members %s: %s", name, err)
		}
		schemaSetToSlice(users)
	}
	return resourceSlackUserGroupRead(ctx, d, m)
}

func resourceSlackUserGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*slack.Client)

	id := d.Id()
	_, err := client.DisableUserGroupContext(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateUserGroupData(d *schema.ResourceData, userGroup slack.UserGroup) diag.Diagnostics {
	if userGroup.ID == "" {
		return diag.Errorf("error setting id: returned usergroup does not have an id")
	}
	d.SetId(userGroup.ID)

	if err := d.Set("name", userGroup.Name); err != nil {
		return diag.Errorf("error setting name: %s", err)
	}

	if err := d.Set("handle", userGroup.Handle); err != nil {
		return diag.Errorf("error setting handle: %s", err)
	}

	if err := d.Set("description", userGroup.Description); err != nil {
		return diag.Errorf("error setting description: %s", err)
	}

	if err := d.Set("channels", userGroup.Prefs.Channels); err != nil {
		return diag.Errorf("error setting channels: %s", err)
	}

	if err := d.Set("users", userGroup.Users); err != nil {
		return diag.Errorf("error setting users: %s", err)
	}

	return nil
}
