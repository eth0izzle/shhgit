rule java_keystore_file
{
	meta:
		author = "Paul Price"
		description = "Finds Java keystore file"

	condition:
		extension contains ".jks"
}