rule postgresql_client_command_history_file
{
	meta:
		author = "Paul Price"
		description = "Finds PostgreSQL client command history file"

	condition:
		filename matches /^\.?psql_history$/
}