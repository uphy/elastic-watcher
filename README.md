# Elastic Watcher

Elastic Watcher watches the Elasticsearch for the changes and perform some actions.  
The watch config file format based on [Elastic X-Pack watcher feature](https://www.elastic.co/guide/en/x-pack/current/xpack-alerting.html)

## Defferences

- The default script language is not `painless` but `javascript`.
- Payload split
- Input condition
- dry_run option for action
- TO/CC email address can be JSON array format.