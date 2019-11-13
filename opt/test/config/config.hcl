
# Core configuration
core {
	address 				= "0.0.0.0"
	port      			= "443"

	log_level  			= "3"
	log_out    			= "stdout"

	file_directory	= "./"

	# Metrics settings
	metrics {
		enabled = true
		path = "/MicroMetrics"
		out = "./metrics.json"
		users {
			"Admin" = "test123"
		}
	}

	virtual_hosting 	= true
	virtual_hosts {
		"localhost" = "./opt/test/config/vhost1.hcl"
		"127.0.0.1" = "./opt/test/config/vhost2.hcl"
	}

	tls {
		enabled 				= true
		tls_cert				= "./opt/test/certs/localhost.cert"
		tls_key					= "./opt/test/certs/localhost.key"
	}
}
