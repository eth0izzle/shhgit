rule rubygems_credentials_file
{
	meta:
		author = "Paul Price"
		description = "Finds Rubygems credentials file"

	condition:
		filepath matches /\.?gem\/credentials$/
}