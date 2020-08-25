rule esmtp_configuration
{
	meta:
		author = "Paul Price"
		description = "Finds esmtp configuration"

	condition:
		filename matches /.esmtprc$/
}