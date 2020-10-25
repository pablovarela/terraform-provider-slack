data slack_user test_user_00 {
  name = "contact_test-user-ter"
}

data slack_user test_user_01 {
  name = "contact_test-user-206"
}

resource slack_conversation aws_chatbot {
  name              = "aws-chat-bot-notifications-pablo"
  topic             = "AWS ChatBot Notifications"
  permanent_members = [data.slack_user.test_user_01.id, data.slack_user.test_user_00.id]
  is_private        = true
}
