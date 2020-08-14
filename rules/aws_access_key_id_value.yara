rule aws_access_key_id_value
{
	meta:
		author = "Paul Price"
		description = "Finds AWS Access Key ID Value"

	strings:
		$ = /(A3T[A-Z0-9]|AKIA|AGPA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}/

	condition:
		any of them
}