# Surf Automation

It's web service which provides API and should be integrated with Slack.

This service provides some Slack Shorcuts, Slack Slack Commands and other automations which is very useful for IT company. 

## Commands

### CTO Feedback

Or just a feedback in common. It's useful for gathering anonymous feedback about one specific person. In my case - about CTO.

How does it works:
1. Handles `POST /slack/cto/feedback` utl path (Slack Slach Command call it)
2. Reads payload from slack's web hook (only text which was sent to command)
3. And then sends feedback text to another user

To determine the user who will recive feedback (CTO) use `SLACK_CHANNEL_ID_FOR_NOTIFICATIONS`

## Shortcuts

### Confusing!

We use it for reporting some confusing requests to system administrators.

For example when people say nothing about details, etc. 

How does it works:
1. An SA runs shortcut in slack (the shortcut should be runned for specific message)
2. Slack calls `POST /commands/ops/wtf`
3. Service validates the reporting:
    1. Report was sent from specific channel (which is listed in `ALLOWED_CHANNELS_IDS`)
    2. Report was sent by specific person (whose id is listed in `ALLOWED_REPORTERS_IDS`)
    3. Author of reported message is reportable (author's id is **not** listed in `OPS_WTF_RESTRICTED_AUTHORS_IDS`)
4. Service save the report to MongoDB which is connected by `MONGODB_CONNECTION_STRING`
5. Service sents information about report (just about event) to specific channel (which id was written in `SLACK_CHANNEL_ID_FOR_NOTIFICATIONS`)

## Automatizations

### 0-level support

**WARNING** 

Don't add steve to any random channel, because this automatization will push messages in any channel where `Steve` is

1. Listen for all channels in which bot `Steve` was added
2. If new message was posted to one of this channel the `POST /slack_events/steve` will be called by slack
3. Service validate that author of this message wasn't this bot (vai comparing author id with this bot id which was set in `STEVE_SLACK_BOT_ID`)
4. Validates that this message wasn't in thread
5. Send specific message in thread linked to the user's message (text fo the message should be set in `STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY`)

Slack api auth token must be set in `STEVE_SLACK_BOT_AUTH_TOKEN`

---

All `ENV` parameters:

`STEVE_SLACK_BOT_ID` - Id of the bot. This id should be the id of user whose token you use in `STEVE_SLACK_BOT_AUTH_TOKEN`

`STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY` - Preformatted message which will be sent for any new message in shannels where `Steve` is

`STEVE_SLACK_BOT_AUTH_TOKEN` - Slack API Auth token. Should belong to user which id you use in `STEVE_SLACK_BOT_ID`

`OPS_WTF_RESTRICTED_AUTHORS_IDS` - Id of authors whe can't be reported by `Confusing!` (elements should be splitted by `,`)

`ALLOWED_REPORTERS_IDS` - Idsof users who can report messages via `Confusing!` (elements should be splitted by `,`)

`ALLOWED_CHANNELS_IDS` - Ids of channels from which messages can be reported via `Confusing!` (elements should be splitted by `,`)

`MONGODB_CONNECTION_STRING` - Connection string with credentials to MongoDB instance

`SLACK_CHANNEL_ID_FOR_NOTIFICATIONS` - channel (or user) id who should recieve information about success events or feedback