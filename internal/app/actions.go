package app

import (
	"log"
	"time"

	"github.com/muonsoft/validation/validate"
	"github.com/omega1x/crta-box/internal/toolkit"
	"github.com/urfave/cli/v2"
)

func actionPingBox(cCtx *cli.Context) error {

	var (
		host     string = cCtx.String("ipv4")
		port     uint16 = uint16(cCtx.Uint("port"))
		database string = cCtx.String("database")
		username string = cCtx.String("username")
		password string = cCtx.String("password")
	)

	if err := validate.IPv4(host); err != nil {
		log.Fatalf(LOG_ERROR+"[%s] is invalid IPv4-address", host)
	}

	if port != BOX_DEFAULT_PORT {
		log.Printf(LOG_WARN+"[%d] is non-standard port has been chosen for *CRTA-BOX* server", port)
	}

	log.Printf(LOG_INFO+"Check connection to *CRTA-BOX* on host [%s:%d]", host, port)

	diagnostics, err := toolkit.PingBox(host, port, database, username, password)
	if err != nil {
		log.Fatalf(LOG_ERROR + diagnostics + "; " + err.Error())
	}

	log.Printf(LOG_INFO + "Checked. " + diagnostics)
	return nil
}

func actionPingHouse(cCtx *cli.Context) error {
	var (
		host     string = cCtx.String("ipv4")
		port     uint16 = uint16(cCtx.Uint("port"))
		database string = cCtx.String("database")
		table    string = cCtx.String("table")
		username string = cCtx.String("username")
		password string = cCtx.String("password")
	)

	if err := validate.IPv4(host); err != nil {
		log.Fatalf(LOG_ERROR+"[%s] is invalid IPv4-address", host)
	}

	if port != HOUSE_DEFAULT_PORT {
		log.Printf(LOG_WARN+"[%d] is non-standard port for *ClickHouse* database", port)
	}

	log.Printf(LOG_INFO+"Check connection to *ClickHouse* on host [%s:%d]", host, port)

	diagnostics, err := toolkit.PingHouse(host, port, database, table, username, password)
	if err != nil {
		log.Fatalf(LOG_ERROR + diagnostics + "; " + err.Error())
	}

	log.Printf(LOG_INFO + "Checked. " + diagnostics)
	return nil
}

// Data stream action from *CRTA-BOX* to *ClickHouse*-database
func actionStream(cCtx *cli.Context) error {
	var (
		box_host     string = cCtx.String("box_ipv4")
		box_port     uint16 = uint16(cCtx.Uint("box_port"))
		box_database string = cCtx.String("box_database")
		box_username string = cCtx.String("box_username")
		box_password string = cCtx.String("box_password")

		house_host     string = cCtx.String("house_ipv4")
		house_port     uint16 = uint16(cCtx.Uint("house_port"))
		house_database string = cCtx.String("house_database")
		house_table    string = cCtx.String("house_table")
		house_username string = cCtx.String("house_username")
		house_password string = cCtx.String("house_password")
	)

	// Check connection data
	//  * box
	if err := validate.IPv4(box_host); err != nil {
		log.Fatalf(LOG_ERROR+"[%s] is invalid IPv4-address", box_host)
	}

	if box_port != BOX_DEFAULT_PORT {
		log.Printf(LOG_WARN+"[%d] is non-standard port for *CRTA-BOX* server", box_port)
	}

	// * house
	if err := validate.IPv4(house_host); err != nil {
		log.Fatalf(LOG_ERROR+"[%s] is not valid IPv4-address.", house_host)
	}

	if house_port != HOUSE_DEFAULT_PORT {
		log.Printf(LOG_WARN+"[%d] is non-standard port for *ClickHouse*-database", house_port)
	}

	// Enter streaming process
	log.Printf(LOG_INFO+"Stream data from *CRTA-BOX* on host [%s:%d] to *ClickHouse* on host [%s:%d]", box_host, box_port, house_host, house_port)

	var (
		JobCount     uint64 = 0
		FailJobCount uint64 = 0
	)

	for FailJobCount < JOB_MAX_FAIL_COUNT {
		JobCount++
		log.Printf(LOG_INFO+"Do stream job [%d]", JobCount)

		house_last_record, transmit_count_row, transmit_count_row_fail_Scan, transmit_count_row_fail_OutOfRange, transmit_count_row_fail_Batch,
			transmit_timestamp_first, transmit_timestamp_last, err := toolkit.Transmit(
			box_host,
			box_port,
			box_database,
			box_username,
			box_password,
			house_host,
			house_port,
			house_database,
			house_table,
			house_username,
			house_password,
		)

		if err != nil {
			log.Printf(LOG_ERROR+"Job [%d] failed with message: %v", JobCount, err)
			log.Printf(LOG_WARN+"All previous [%d] job(s) already failed too!", FailJobCount)
			FailJobCount++
			continue
		}
		FailJobCount = 0
		log.Printf(
			LOG_INFO+"Job [%d]: *ClickHouse*-database last record timestamp is [%d]", JobCount, house_last_record,
		)
		log.Printf(
			LOG_INFO+"Job [%d]: [%d] rows for timestamps between [%d - %d] transmitted",
			JobCount,
			transmit_count_row,
			transmit_timestamp_first,
			transmit_timestamp_last,
		)
		if transmit_count_row_fail_Scan > 0 {
			log.Printf(
				LOG_WARN+"Job [%d]: [%d]/[%d] failed scans",
				JobCount,
				transmit_count_row_fail_Scan,
				transmit_count_row,
			)
		}
		if transmit_count_row_fail_OutOfRange > 0 {
			log.Printf(
				LOG_WARN+"Job [%d]: [%d]/[%d] out of range",
				JobCount,
				transmit_count_row_fail_OutOfRange,
				transmit_count_row,
			)
		}
		if transmit_count_row_fail_Batch > 0 {
			log.Printf(
				LOG_WARN+"Job [%d]: [%d]/[%d] failed batch",
				JobCount,
				transmit_count_row_fail_Batch,
				transmit_count_row,
			)
		}

		if transmit_count_row_fail_Scan+transmit_count_row_fail_OutOfRange+transmit_count_row_fail_Batch == 0 {
			log.Printf(
				LOG_INFO+"Ok. Job [%d]: [%d] rows for period [%s <-> %s] successfully transmitted",
				JobCount,
				transmit_count_row,
				time.Unix(transmit_timestamp_first, 0).Format(time.RFC1123),
				time.Unix(transmit_timestamp_last, 0).Format(time.RFC1123),
			)
		}

		log.Printf(LOG_INFO+"Now sleep for [%d] seconds", JOB_SLEEP_DURATION)
		time.Sleep(JOB_SLEEP_DURATION * time.Second)
	}

	log.Printf(LOG_WARN + "Too many failed jobs. No sense to continue")
	return nil
}
