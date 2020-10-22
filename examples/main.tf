provider slack {
  token   = var.slack_token
}

data slack_user aws_chat_bot {
  name = "aws"
}

resource slack_conversation aws_chatbot {
  name       = "AWS Chatbot Notifications"
  topic      = "AWS ChatBot Notifications"
  members    = [data.slack_user.aws_chat_bot.id]
  is_private = true
}