rule terraform_variable_config_file
{
	meta:
		author = "Paul Price"
		description = "Finds Terraform variable config file"

	condition:
		filename contains "terraform.tfvars"
}