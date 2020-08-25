rule recon_ng_web_reconnaissance_framework_api_key_database
{
	meta:
		author = "Paul Price"
		description = "Finds Recon-ng web reconnaissance framework API key database"

	condition:
		filepath matches /\.?recon-ng\/keys\.db$/
}