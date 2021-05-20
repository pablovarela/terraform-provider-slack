data "slack_user" "test_user_00" {
  name = "contact_test-user-ter"
}

data "slack_user" "test_user_01" {
  name = "contact_test-user-206"
}

resource "slack_usergroup" "my_group" {
  name        = "TestGroup"
  handle      = "test"
  description = "Test user group"
  users       = [data.slack_user.test_user_01.id, data.slack_user.test_user_00.id]
}

resource "slack_conversation" "my_conversation" {
  name            = "my-conversation"
  topic           = "some interesting topic"
  private_members = slack_usergroup.my_group.users
  is_private      = true
}

data "slack_usergroup" "by_name" {
  name = slack_usergroup.my_group.name
}

data "slack_usergroup" "by_id" {
  usergroup_id = slack_usergroup.my_group.id
}
