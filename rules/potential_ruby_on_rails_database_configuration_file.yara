rule potential_ruby_on_rails_database_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Potential Ruby On Rails database configuration file"

	condition:
		filename contains "database.yml"
}