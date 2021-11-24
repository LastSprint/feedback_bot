# Surf Automation
sdfsdfsdf
----
----
BBBBBB
----
----

It's web service which provides API and should be integrated with Slack.

This service provides some Slack Shorcuts, Slack Slack Commands and other automations which is very useful for IT company. 

## Commands

| Command | Input Args| Description |   
|----------|---------|-----|
|`/cto_feedback`| `/cto_feedback and just text...` |Send feedback about Surf CTO to him |
|`/ops_and_sa_weekly`| - | Send SA and DevOps work digest into chat |
|`/steve_analyze_wl_btw_projects`| `/steve_analyze_wl_btw_projects jiraLogin1 jiraLogin2...` | Result: time that person spent for last 7 days distributed between his projects |
| `/steve_tech_ws` | - | Print internal technologies adaptation|
### CTO Feedback

Command: `/cto_feedback`

Input: just a feedback text after command

Just a feedback in common. It's useful for gathering anonymous feedback about one specific person. In my case - about CTO.

How does it works:
1. Handles `POST /slack/cto/feedback` utl path (Slack Slach Command call it)
2. Reads payload from slack's web hook (only text which was sent to command)
3. And then sends feedback text to another user

To determine the user who will recive feedback (CTO) use `SLACK_CHANNEL_ID_FOR_NOTIFICATIONS`

### Ops And SA Weekly

Command: `/ops_and_sa_weekly`

Reports:
- How many user' requests were
- What kind of reports were and how many

Use `STEVE_DEVOPS_AND_SA_CHANNEL_ID` to choose a channel from which channel reports and requests shout be counted

### Analyze WorkLog For Projects

Command: `/steve_analyze_wl_btw_projects`

Input: Jira logins of users whose work log you want to see (logins splitteed by whitespace)

Result: time that person spent for last 7 days in form:

```
$USERNAME$ spent:
- $PROJECT$: $TIME$
...
```

### Technologies Adaptation

Command: `/steve_tech_ws`

Will print adaptation of `Templates` and `SurfGen`

## Shortcuts

### Confusing!

CallbackID: `ops_wtf`

ReportType: `bad_request`

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

### JenkinsLogsWasNotRead

CallbackID: `ops_didnt_read_jenkins_log`

ReportType: `did_not_read_jenkins_logs`

Should be called if somebody create request about errors in Jenkins' build, but he/she didn't event open logs page.

The algorithm is the same as in `Confusign!` (actually the code is the same)

## Automatizations

### 0-level support

1. Listen for all channels in which bot `Steve` was added
2. If new message was posted to one of this channel the `POST /slack_events/steve` will be called by slack
3. Service validate that author of this message wasn't this bot (via comparing author id with this bot id which was set in `STEVE_SLACK_BOT_ID`)
4. Validates that this message wasn't in thread
5. Validates that message was sent from channel listed in `SUPPORT_AUTOMATION_CHANNELS_TO_REPLY`
5. Send specific message in thread linked to the user's message (text fo the message should be set in `STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY`)
6. Increment requests count in MongoDB

---

Slack api auth token must be set in `STEVE_SLACK_BOT_AUTH_TOKEN`

---

All `ENV` parameters:


| ENV  key |  Description |   
|----------|--------------|
|`STEVE_SLACK_BOT_ID`|Id of the bot. This id should be the id of user whose token you use in `STEVE_SLACK_BOT_AUTH_TOKEN`|
|`STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY`|Preformatted message which will be sent for any new message in shannels where `Steve` is|
|`SUPPORT_AUTOMATION_CHANNELS_TO_REPLY`|Array of channel ids in which `Steve` can reply on event of creating message|
|`STEVE_SLACK_BOT_AUTH_TOKEN`|Slack API Auth token. Should belong to user which id you use in `STEVE_SLACK_BOT_ID`|
|`OPS_WTF_RESTRICTED_AUTHORS_IDS`|Id of authors whe can't be reported by `Confusing!` (elements should be splitted by `,`)|
|`ALLOWED_REPORTERS_IDS`|Ids of users who can report messages via `Confusing!` (elements should be splitted by `,`)|
|`ALLOWED_CHANNELS_IDS`|Ids of channels from which messages can be reported via `Confusing!` (elements should be splitted by `,`)|
|`MONGODB_CONNECTION_STRING`|Connection string with credentials to MongoDB instance|
|`SLACK_CHANNEL_ID_FOR_NOTIFICATIONS`|Channel (or user) id who should recieve information about success events or feedback|
|`STEVE_DEVOPS_AND_SA_CHANNEL_ID`| used by `Ops And SA Weekly` to choose a channel from which channel reports and requests shout be counted |
|`JIRA_AUTH_API_TOKEN`| Used to get access to JIRA API |
