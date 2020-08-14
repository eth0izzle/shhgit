rule password_safe_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds Password Safe database file"

	condition:
		extension contains ".psafe3"
}