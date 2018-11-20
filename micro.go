package main

import "github.com/redmaner/smux"

// Micro struct which holds all the information of the server
// The micro struct has it's own functions that have access to this data
type micro struct {
	config microConfig
	vhosts map[string]microConfig
	md     metricsData
	router *smux.SRouter
}
