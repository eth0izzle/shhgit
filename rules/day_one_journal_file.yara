rule day_one_journal_file
{
	meta:
		author = "Paul Price"
		description = "Finds Day One journal file"

	condition:
		extension contains ".dayone"
}