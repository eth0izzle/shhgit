rule rails_master_key_used_for_decrypting_credentials_yml_enc_for_rails_5_2_
{
	meta:
		author = "Paul Price"
		description = "Finds Rails master key (used for decrypting credentials.yml.enc for Rails 5.2+)"

	condition:
		filepath matches /ruby[\\\/]config[\\\/]master.key$/
}