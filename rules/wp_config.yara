rule wp_config
{
	meta:
		author = "Paul Price"
		description = "Finds WP-Config"

	strings:
		$ = /define(.{0,20})?(DB_CHARSET|NONCE_SALT|LOGGED_IN_SALT|AUTH_SALT|NONCE_KEY|DB_HOST|DB_PASSWORD|AUTH_KEY|SECURE_AUTH_KEY|LOGGED_IN_KEY|DB_NAME|DB_USER)(.{0,20})?['|"].{10,120}['|"]/

	condition:
		any of them
}