rule sqlite_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds SQLite database file"

	condition:
		extension contains ".sqlite"
}