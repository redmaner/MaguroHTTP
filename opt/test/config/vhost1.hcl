
# Proxy settings
proxy {
	enabled 				= true
	rules {
		"localhost" = "https://www.edsn.nl"
	}
	methods = ["GET", "POST"]
}

# Guard settings
guard {
	rate 						= 100
	rate_burst				= 10
}
