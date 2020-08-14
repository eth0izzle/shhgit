rule potential_puttygen_private_key
{
	meta:
		author = "Paul Price"
		description = "Finds Potential PuTTYgen private key"

	condition:
		extension contains ".ppk"
}