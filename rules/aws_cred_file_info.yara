rule aws_cred_file_info
{
	meta:
		author = "Paul Price"
		description = "Finds AWS cred file info"

	strings:
		$ = /(aws_access_key_id|aws_secret_access_key)(.{0,20})?=.[0-9a-zA-Z\/+]{20,40}/

	condition:
		any of them
}