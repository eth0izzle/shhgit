rule potential_linux_passwd_file
{
	meta:
		author = "Paul Price"
		description = "Finds Potential Linux passwd file"

	condition:
		filepath matches /etc\/passwd$/
}