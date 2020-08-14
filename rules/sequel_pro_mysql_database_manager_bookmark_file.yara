rule sequel_pro_mysql_database_manager_bookmark_file
{
	meta:
		author = "Paul Price"
		description = "Finds Sequel Pro MySQL database manager bookmark file"

	condition:
		filename contains "Favorites.plist"
}