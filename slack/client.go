package slack

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/slack-go/slack"
)

// SlackClient wraps the slack.Client with caching for usergroups
// to avoid hitting rate limits when fetching the same data multiple times
// during a single plan/apply cycle.
type SlackClient struct {
	*slack.Client

	userGroupsMu    sync.RWMutex
	userGroupsCache map[userGroupsCacheKey][]slack.UserGroup
}

type userGroupsCacheKey struct {
	includeDisabled bool
	includeUsers    bool
}

// NewSlackClient creates a new SlackClient with caching enabled
func NewSlackClient(token string) *SlackClient {
	return &SlackClient{
		Client:          slack.New(token),
		userGroupsCache: make(map[userGroupsCacheKey][]slack.UserGroup),
	}
}

// GetUserGroupsCached returns usergroups from cache if available,
// otherwise fetches from the API and caches the result.
// The cache is keyed by the includeDisabled and includeUsers options.
func (c *SlackClient) GetUserGroupsCached(ctx context.Context, includeDisabled, includeUsers bool) ([]slack.UserGroup, error) {
	key := userGroupsCacheKey{
		includeDisabled: includeDisabled,
		includeUsers:    includeUsers,
	}

	// Check cache first
	c.userGroupsMu.RLock()
	if cached, ok := c.userGroupsCache[key]; ok {
		c.userGroupsMu.RUnlock()
		return cached, nil
	}
	c.userGroupsMu.RUnlock()

	// Fetch from API with rate limit retry
	opts := []slack.GetUserGroupsOption{
		slack.GetUserGroupsOptionIncludeDisabled(includeDisabled),
		slack.GetUserGroupsOptionIncludeUsers(includeUsers),
	}

	var userGroups []slack.UserGroup
	for {
		var err error
		userGroups, err = c.Client.GetUserGroupsContext(ctx, opts...)
		if err == nil {
			break
		}

		if rateLimitErr, ok := err.(*slack.RateLimitedError); ok {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled while waiting for rate limit: %w", ctx.Err())
			case <-time.After(rateLimitErr.RetryAfter):
				continue
			}
		}
		return nil, err
	}

	// Cache the result
	c.userGroupsMu.Lock()
	c.userGroupsCache[key] = userGroups
	c.userGroupsMu.Unlock()

	return userGroups, nil
}

// InvalidateUserGroupsCache clears the usergroups cache.
// This should be called after any mutation (create, update, delete)
// to ensure subsequent reads get fresh data.
func (c *SlackClient) InvalidateUserGroupsCache() {
	c.userGroupsMu.Lock()
	c.userGroupsCache = make(map[userGroupsCacheKey][]slack.UserGroup)
	c.userGroupsMu.Unlock()
}
