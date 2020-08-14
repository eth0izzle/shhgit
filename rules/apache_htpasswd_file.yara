rule apache_htpasswd_file
{
	meta:
		author = "Paul Price"
		description = "Finds Apache htpasswd file"

	condition:
		filename matches /^\.?htpasswd$/
}