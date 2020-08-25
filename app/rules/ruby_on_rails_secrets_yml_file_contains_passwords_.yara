rule ruby_on_rails_secrets_yml_file_contains_passwords_
{
	meta:
		author = "Paul Price"
		description = "Finds Ruby on rails secrets.yml file (contains passwords)"

	condition:
		filepath matches /web[\\\/]ruby[\\\/]secrets.yml/
}