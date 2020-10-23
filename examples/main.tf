provider slack {
  token   = var.slack_token
}

data slack_user aws_chat_bot {
  name = "aws"
}

resource slack_conversation aws_chatbot {
  name       = "aws-chat-bot-notifications"
  topic      = "AWS ChatBot Notifications"
  members    = [data.slack_user.aws_chat_bot.id]
  is_private = true
}