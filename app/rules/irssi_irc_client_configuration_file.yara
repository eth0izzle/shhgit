rule irssi_irc_client_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Irssi IRC client configuration file"

	condition:
		filepath matches /\.?irssi\/config$/
}