rule carrierwave_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Carrierwave configuration file"

	condition:
		filename contains "carrierwave.rb"
}