rule robomongo_mongodb_manager_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Robomongo MongoDB manager configuration file"

	condition:
		filename contains "robomongo.json"
}