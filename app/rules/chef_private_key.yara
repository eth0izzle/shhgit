rule chef_private_key
{
	meta:
		author = "Paul Price"
		description = "Finds Chef private key"

	condition:
		filepath matches /\.?chef\/(.*)\.pem$/
}