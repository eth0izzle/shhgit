rule shell_profile_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Shell profile configuration file"

	condition:
		filename matches /^\.?(bash_|zsh_)?profile$/
}