rule apple_keychain_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds Apple Keychain database file"

	condition:
		extension contains ".keychain"
}