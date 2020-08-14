rule gnucash_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds GnuCash database file"

	condition:
		extension contains ".gnucash"
}