/*
Package app wraps functions of the `crta-box` console application as a command line
utility, and provides functionality for interacting with the system and the user.
*/
package app

const (
	BOX_DEFAULT_HOST     = "127.0.0.1"
	BOX_DEFAULT_PORT     = 5432
	BOX_DEFAULT_DATABASE = "logger"
	BOX_DEFAULT_USERNAME = "user"
	BOX_DEFAULT_PASSWORD = "pass"

	HOUSE_DEFAULT_HOST     = "127.0.0.1"
	HOUSE_DEFAULT_PORT     = 9000
	HOUSE_DEFAULT_DATABASE = "BOXes"
	HOUSE_DEFAULT_TABLE    = "box00"
	HOUSE_DEFAULT_USERNAME = "user"
	HOUSE_DEFAULT_PASSWORD = "pass"

	LOG_INFO  = "[INFO] "
	LOG_WARN  = "[WARN] "
	LOG_ERROR = "[ERROR] "

	JOB_MAX_FAIL_COUNT = 5
	JOB_SLEEP_DURATION = 250 // [seconds]
)
