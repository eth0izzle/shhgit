rule dbeaver_sql_database_manager_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds DBeaver SQL database manager configuration file"

	condition:
		filename matches /^\.?dbeaver-data-sources.xml$/
}