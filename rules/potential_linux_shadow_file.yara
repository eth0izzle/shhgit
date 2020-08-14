rule potential_linux_shadow_file
{
	meta:
		author = "Paul Price"
		description = "Finds Potential Linux shadow file"

	condition:
		filepath matches /etc\/shadow$/
}