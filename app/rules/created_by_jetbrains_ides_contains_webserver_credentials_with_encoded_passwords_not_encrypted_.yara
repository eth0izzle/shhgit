rule created_by_jetbrains_ides_contains_webserver_credentials_with_encoded_passwords_not_encrypted_
{
	meta:
		author = "Paul Price"
		description = "Finds Created by Jetbrains IDEs, contains webserver credentials with encoded passwords (not encrypted!)"

	condition:
		filepath matches /\.?idea[\\\/]WebServers.xml$/
}