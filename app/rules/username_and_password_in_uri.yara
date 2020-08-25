rule username_and_password_in_uri
{
	meta:
		author = "Paul Price"
		description = "Finds Username and password in URI"

	strings:
		$ = /([\w+]{3,24})(:\/\/)([^$<]{1})([^\s";]{1,}):([^$<]{1})([^\s";\/]{1,})@[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,24}([^\s]+)/

	condition:
		any of them
}