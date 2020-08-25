rule microsoft_bitlocker_trusted_platform_module_password_file
{
	meta:
		author = "Paul Price"
		description = "Finds Microsoft BitLocker Trusted Platform Module password file"

	condition:
		extension contains ".tpm"
}