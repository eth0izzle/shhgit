rule kde_wallet_manager_database_file
{
	meta:
		author = "Paul Price"
		description = "Finds KDE Wallet Manager database file"

	condition:
		extension contains ".kwallet"
}