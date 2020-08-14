rule heroku_config_file
{
	meta:
		author = "Paul Price"
		description = "Finds Heroku config file"

	condition:
		filename contains "heroku.json"
}