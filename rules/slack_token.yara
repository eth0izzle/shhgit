rule slack_token
{
	meta:
		author = "Paul Price"
		description = "Finds Slack Token"

	strings:
		$ = /(xox[pboa]\-[0-9]{12}\-[0-9]{12}\-[0-9]{12}\-[a-z0-9]{32})/

	condition:
		any of them
}