rule potential_cryptographic_key_bundle
{
	meta:
		author = "Paul Price"
		description = "Finds Potential cryptographic key bundle"

	condition:
		extension contains ".asc"
}