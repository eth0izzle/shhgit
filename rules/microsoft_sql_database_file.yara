rule microsoft_sql_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds Microsoft SQL database file"

	condition:
		extension contains ".mdf"
}