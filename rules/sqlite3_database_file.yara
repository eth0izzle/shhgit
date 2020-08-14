rule sqlite3_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds SQLite3 database file"

	condition:
		extension contains ".sqlite3"
}