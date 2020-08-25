rule windows_bitlocker_full_volume_encrypted_data_file
{
	meta:
		author = "Paul Price"
		description = "Finds Windows BitLocker full volume encrypted data file"

	condition:
		extension contains ".fve"
}