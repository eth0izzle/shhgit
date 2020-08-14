rule potential_jrnl_journal_file
{
	meta:
		author = "Paul Price"
		description = "Finds Potential jrnl journal file"

	condition:
		filename contains "journal.txt"
}