rule facebook_client_id
{
	meta:
		author = "Paul Price"
		description = "Finds Facebook Client ID"

	strings:
		$ = /(facebook|fb)(.{0,20})?['\"][0-9]{13,17}['\"]/

	condition:
		any of them
}