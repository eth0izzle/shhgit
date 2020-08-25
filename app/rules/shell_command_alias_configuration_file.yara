rule shell_command_alias_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Shell command alias configuration file"

	condition:
		filename matches /^\.?(bash_|zsh_)?aliases$/
}