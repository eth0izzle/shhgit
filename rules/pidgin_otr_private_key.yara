rule pidgin_otr_private_key
{
	meta:
		author = "Paul Price"
		description = "Finds Pidgin OTR private key"

	condition:
		filename contains "otr.private_key"
}