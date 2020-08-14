rule public_ssh_key
{
	meta:
		author = "Paul Price"
		description = "Finds Public ssh key"

	condition:
		filename contains "id_rsa_pub"
}