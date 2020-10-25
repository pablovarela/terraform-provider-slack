package slack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/slack-go/slack"
	"os"
	"testing"
)

var slackClient *slack.Client

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedSlackClient() (interface{}, error) {
	if slackClient != nil {
		return slackClient, nil
	}

	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		return nil, fmt.Errorf("could not initialize Slack client. Set environment variable SLACK_TOKEN")
	}

	return slack.New(token), nil
}
