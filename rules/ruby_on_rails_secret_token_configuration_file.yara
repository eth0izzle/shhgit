rule ruby_on_rails_secret_token_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Ruby On Rails secret token configuration file"

	condition:
		filename contains "secret_token.rb"
}