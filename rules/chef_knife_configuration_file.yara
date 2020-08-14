rule chef_knife_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Chef Knife configuration file"

	condition:
		filename contains "knife.rb"
}