rule log_file
{
	meta:
		author = "Paul Price"
		description = "Finds Log file"

	condition:
		extension contains ".log"
}