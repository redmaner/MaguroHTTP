### r3
* new: server side read and write timeouts
* new: metrics can only be used with TLS Enabled
* fixed: MIME types have been improved to better match web standards
* fixed: use ServeDir also for requests for directories
* fixed: graceful stop of server when not using TLS

### r2
* new: Support for automated certificates with Let's Encrypt
* new: restructured, more logical, configuration: breaks compatibility with r1 but is more suited for the future
* fixed: proxy not setting proper headers which could break certain websites
* fixed: flushing MicroMetrics while MicroMetrics wasn't enabled
* fixed: incompliance with golint

### r1
* initial release
