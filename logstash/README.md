# Logstash Configuration

## 1. Install ElasticSearch, Kibana, and Logstash
- Follow the instructions on Elastic's website: https://www.elastic.co/guide/en/elastic-stack-get-started/current/get-started-elastic-stack.html

## 2. Configure Logstash

In the folder logstash/config/01-shhgit-pipeline.conf is a working example with a few Grok rules to process the incoming messages from shhgit.  These Grok rules are imperfect, but should provide a starting point to improve upon.  

- 01-shhgit-pipeline.conf should be placed in the Logstash /etc/logstash/conf.d directory.
- Update the output elasticsearch to point to your elasticsearch cluster

## 3. Configure Kibana

### Configure Pattern
- Click the gear icon (management) in the lower left
- Click Kibana -> Index Patters
- Click Create New Index Pattern
- Type "shhgit-*" into the input box, then click Next Step

### Import Dashboard
In your web browser go to Kibana's IP using port 5601 (ex: 192.168.0.1:5601)
- Click Management -> Saved Objects
- You can import the dashboard found in the dashboard folder via the Import button in the top-right corner.
  - shhgit Dashboard
  
## TODO
- Clean up Grok Rules
- Add additional Grok Rules
