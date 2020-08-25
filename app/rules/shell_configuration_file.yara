rule shell_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Shell configuration file"

	condition:
		filename matches /^\.?(bash|zsh|csh)rc$/
}