rule keepass_password_manager_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds KeePass password manager database file"

	condition:
		extension matches /^kdbx?$/
}