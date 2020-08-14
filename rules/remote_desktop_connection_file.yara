rule remote_desktop_connection_file
{
	meta:
		author = "Paul Price"
		description = "Finds Remote Desktop connection file"

	condition:
		extension contains ".rdp"
}