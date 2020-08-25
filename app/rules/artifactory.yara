rule artifactory
{
	meta:
		author = "Paul Price"
		description = "Finds Artifactory"

	strings:
		$ = /artifactory.{0,50}(\"|'|`)?[a-zA-Z0-9=]{112}(\"|'|`)?/

	condition:
		any of them
}