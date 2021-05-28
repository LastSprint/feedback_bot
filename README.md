# feedback_bot

It's just a very simple service which is used to handle Slack commands.
Service just get `text` from slack's command payload and save it into a file.

You can specify which file you prefer via env variable `FEEDBACK_BOT_DB_FILE_PATH`

Or you can use `Dockerfile` from this repo
