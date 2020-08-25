rule mysql_client_command_history_file
{
	meta:
		author = "Paul Price"
		description = "Finds MySQL client command history file"

	condition:
		filename matches /^\.?mysql_history$/
}