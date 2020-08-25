rule sql_dump_file
{
	meta:
		author = "Paul Price"
		description = "Finds SQL dump file"

	condition:
		extension matches /^sql(dump)?$/
}