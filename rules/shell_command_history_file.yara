rule shell_command_history_file
{
	meta:
		author = "Paul Price"
		description = "Finds Shell command history file"

	condition:
		filename matches /^\.?(bash_|zsh_|sh_|z)?history$/
}