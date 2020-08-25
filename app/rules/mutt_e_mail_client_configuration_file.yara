rule mutt_e_mail_client_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Mutt e-mail client configuration file"

	condition:
		filename matches /^\.?muttrc$/
}