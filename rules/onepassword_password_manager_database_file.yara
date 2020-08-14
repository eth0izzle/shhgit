rule onepassword_password_manager_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds 1Password password manager database file"

	condition:
		extension contains ".agilekeychain"
}