rule ssh_password
{
	meta:
		author = "Paul Price"
		description = "Finds SSH Password"

	strings:
		$ = /sshpass -p.*['|\"]/

	condition:
		any of them
}