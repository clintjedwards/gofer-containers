package main

import (
	"log"
)

// Generally notifiers have several env vars that are passed in
// 1) The env vars of the main config, which is set in the Gofer config file.
// 	ex. "GOFER_NOTIFIER_MAIN_CONFIG_<VALUE>"
// 2) The env vars of the user supplied config, which is set by users in the pipeline config file.
// 	ex. "GOFER_NOTIFIER_USER_CONFIG_<VALUE>"
// 3) The "GOFER_API_TOKEN" which is passed in so that the notifier can request more information about a particular
// run.

func main() {
	log.Printf("notification successful for run")
}
