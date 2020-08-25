rule git_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Git configuration file"

	condition:
		filename matches /^\.?gitconfig$/
}