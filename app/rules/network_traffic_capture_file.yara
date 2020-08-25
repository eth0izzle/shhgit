rule network_traffic_capture_file
{
	meta:
		author = "Paul Price"
		description = "Finds Network traffic capture file"

	condition:
		extension contains ".pcap"
}