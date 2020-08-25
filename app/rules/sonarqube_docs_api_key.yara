rule sonarqube_docs_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds SonarQube Docs API Key"

	strings:
		$ = /sonar.{0,50}(\"|'|`)?[0-9a-f]{40}(\"|'|`)?/

	condition:
		any of them
}