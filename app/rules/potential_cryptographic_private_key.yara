rule potential_cryptographic_private_key
{
	meta:
		author = "Paul Price"
		description = "Finds Potential cryptographic private key"

	condition:
		extension matches /^key(pair)?$/
}