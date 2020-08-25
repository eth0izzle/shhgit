rule configuration_file_for_auto_login_process
{
	meta:
		author = "Paul Price"
		description = "Finds Configuration file for auto-login process"

	condition:
		filename matches /^(\.|_)?netrc$/
}