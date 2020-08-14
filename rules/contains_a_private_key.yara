rule contains_a_private_key
{
	meta:
		author = "Paul Price"
		description = "Finds Contains a private key"

	strings:
		$ = /-----BEGIN (EC|RSA|DSA|OPENSSH|PGP) PRIVATE KEY/

	condition:
		any of them
}