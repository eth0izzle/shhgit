rule hockeyapp
{
	meta:
		author = "Paul Price"
		description = "Finds HockeyApp"

	strings:
		$ = /hockey.{0,50}(\"|'|`)?[0-9a-f]{32}(\"|'|`)?/

	condition:
		any of them
}