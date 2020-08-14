rule private_ssh_key
{
	meta:
		author = "Paul Price"
		description = "Finds Private SSH key"

	condition:
		filename matches /^.*_ecdsa$/
}