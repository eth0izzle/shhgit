rule potential_mediawiki_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Potential MediaWiki configuration file"

	condition:
		filename contains "LocalSettings.php"
}