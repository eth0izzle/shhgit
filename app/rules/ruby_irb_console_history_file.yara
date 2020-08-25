rule ruby_irb_console_history_file
{
	meta:
		author = "Paul Price"
		description = "Finds Ruby IRB console history file"

	condition:
		filename matches /^\.?irb_history$/
}