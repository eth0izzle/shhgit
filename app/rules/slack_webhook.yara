rule slack_webhook
{
	meta:
		author = "Paul Price"
		description = "Finds Slack Webhook"

	strings:
		$ = /https:\/\/hooks.slack.com\/services\/T[a-zA-Z0-9_]{8}\/B[a-zA-Z0-9_]{8}\/[a-zA-Z0-9_]{24}/

	condition:
		any of them
}