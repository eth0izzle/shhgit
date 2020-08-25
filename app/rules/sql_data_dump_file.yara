rule sql_data_dump_file
{
	meta:
		author = "Paul Price"
		description = "Finds SQL Data dump file"

	condition:
		extension contains ".sqldump"
}