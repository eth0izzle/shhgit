rule postgresql_password_file
{
	meta:
		author = "Paul Price"
		description = "Finds PostgreSQL password file"

	condition:
		filename matches /^\.?pgpass$/
}