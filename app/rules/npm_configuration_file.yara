rule npm_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds NPM configuration file"

	condition:
		filename matches /^\.?npmrc$/
}