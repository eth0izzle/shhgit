rule pidgin_chat_client_account_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Pidgin chat client account configuration file"

	condition:
		filepath matches /\.?purple\/accounts\.xml$/
}