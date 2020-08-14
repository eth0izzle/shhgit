rule aws_cli_credentials_file
{
	meta:
		author = "Paul Price"
		description = "Finds AWS CLI credentials file"

	condition:
		filepath matches /\.?aws\/credentials$/
}