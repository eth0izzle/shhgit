rule omniauth_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds OmniAuth configuration file"

	condition:
		filename contains "omniauth.rb"
}