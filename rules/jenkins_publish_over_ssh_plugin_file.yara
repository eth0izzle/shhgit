rule jenkins_publish_over_ssh_plugin_file
{
	meta:
		author = "Paul Price"
		description = "Finds Jenkins publish over SSH plugin file"

	condition:
		filename contains "jenkins.plugins.publish_over_ssh.BapSshPublisherPlugin.xml"
}