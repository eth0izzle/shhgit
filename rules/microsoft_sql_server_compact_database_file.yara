rule microsoft_sql_server_compact_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds Microsoft SQL server compact database file"

	condition:
		extension contains ".sdf"
}