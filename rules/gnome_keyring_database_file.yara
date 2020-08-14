rule gnome_keyring_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds GNOME Keyring database file"

	condition:
		extension matches /^key(store|ring)$/
}