rule hexchat_xchat_irc_client_server_list_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Hexchat/XChat IRC client server list configuration file"

	condition:
		filepath matches /\.?xchat2?\/servlist_?\.conf$/
}