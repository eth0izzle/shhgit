rule firefox_saved_password_collection_can_be_decrypted_using_keys4_db_
{
	meta:
		author = "Paul Price"
		description = "Finds Firefox saved password collection (can be decrypted using keys4.db)"

	condition:
		filepath matches /\.?mozilla[\\\/]firefox[\\\/]logins.json$/
}