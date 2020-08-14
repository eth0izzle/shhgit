rule mongoid_config_file
{
	meta:
		author = "Paul Price"
		description = "Finds Mongoid config file"

	condition:
		filename contains "mongoid.yml"
}