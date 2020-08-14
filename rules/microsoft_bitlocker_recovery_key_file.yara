rule microsoft_bitlocker_recovery_key_file
{
	meta:
		author = "Paul Price"
		description = "Finds Microsoft BitLocker recovery key file"

	condition:
		extension contains ".bek"
}