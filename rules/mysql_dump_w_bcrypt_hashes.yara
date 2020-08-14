rule mysql_dump_w_bcrypt_hashes
{
	meta:
		author = "Paul Price"
		description = "Finds MySQL dump w/ bcrypt hashes"

	condition:
		filename contains "dump.sql"
}