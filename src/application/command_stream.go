package application

import (
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/muonsoft/validation/validate"
	"github.com/omega1x/crta-box/src/com"
	"github.com/omega1x/crta-box/src/conf"
	"github.com/urfave/cli/v2"
)

// Data stream action from *CRTA-BOX* to *ClickHouse*-database
func DoStream(cCtx *cli.Context) error {
	var (
		box_host     string = cCtx.String("box_ipv4")
		box_port     uint16 = uint16(cCtx.Uint("box_port"))
		box_database string = cCtx.String("box_database")
		box_username string = cCtx.String("box_username")
		box_password string = cCtx.String("box_password")

		house_host     string = cCtx.String("house_ipv4")
		house_port     uint16 = uint16(cCtx.Uint("house_port"))
		house_database string = cCtx.String("house_database")
		house_username string = cCtx.String("house_username")
		house_password string = cCtx.String("house_password")
	)

	// Check connection data
	//  * box
	if err := validate.IPv4(box_host); err != nil {
		log.Fatalf(conf.LOG_ERROR+"[%s] is invalid IPv4-address", box_host)
	}

	if box_port != conf.POSTGRES_DEFAULT_PORT {
		log.Printf(conf.LOG_WARN+"[%d] is a non-standard port for *PostgreSQL*", box_port)
	}

	// * house
	if err := validate.IPv4(house_host); err != nil {
		log.Fatalf(conf.LOG_ERROR+"[%s] is not valid IPv4-address.", house_host)
	}

	if house_port != conf.HOUSE_DEFAULT_PORT {
		log.Printf(conf.LOG_WARN+"[%d] is a non-standard port for *ClickHouse*", house_port)
	}

	// Enter streaming process
	log.Printf(conf.LOG_INFO+"Stream data from *CRTA-BOX* on host [%s:%d] to *ClickHouse* on host [%s:%d]...", box_host, box_port, house_host, house_port)

	var (
		JobCount     uint64 = 0
		FailJobCount uint64 = 0
	)

	for FailJobCount < conf.JOB_MAX_FAIL_COUNT {
		JobCount++
		log.Printf(conf.LOG_INFO+"Do job [%d]...", JobCount)

		house_last_record,
			move_count_row,
			move_count_row_fail_Scan,
			move_count_row_fail_OutOfRange,
			move_count_row_fail_Batch,
			move_timestamp_first,
			move_timestamp_last,
			err := com.BoxMove(
			box_host,
			box_port,
			box_database,
			box_username,
			box_password,
			house_host,
			house_port,
			house_database,
			house_username,
			house_password,
		)

		if err != nil {
			log.Printf(conf.LOG_ERROR+"Job [%d] failed with message: %v", JobCount, err)
			log.Printf(conf.LOG_WARN+"All previous [%d] job(s) already failed too!", FailJobCount)
			FailJobCount++
			continue
		}
		FailJobCount = 0
		log.Printf(
			conf.LOG_INFO+"Job [%d]: *ClickHouse* last record timestamp is [%d]", JobCount, house_last_record,
		)
		log.Printf(
			conf.LOG_INFO+"Job [%d]: [%d] rows for timestamps between [%d - %d] moved",
			JobCount,
			move_count_row,
			move_timestamp_first,
			move_timestamp_last,
		)
		if move_count_row_fail_Scan > 0 {
			log.Printf(
				conf.LOG_WARN+"Job [%d]: [%d]/[%d] failed scans",
				JobCount,
				move_count_row_fail_Scan,
				move_count_row,
			)
		}
		if move_count_row_fail_OutOfRange > 0 {
			log.Printf(
				conf.LOG_WARN+"Job [%d]: [%d]/[%d] out of range",
				JobCount,
				move_count_row_fail_OutOfRange,
				move_count_row,
			)
		}
		if move_count_row_fail_Batch > 0 {
			log.Printf(
				conf.LOG_WARN+"Job [%d]: [%d]/[%d] failed batch",
				JobCount,
				move_count_row_fail_Batch,
				move_count_row,
			)
		}

		if move_count_row_fail_Scan+move_count_row_fail_OutOfRange+move_count_row_fail_Batch == 0 {
			log.Printf(
				conf.LOG_INFO+"Ok. Job [%d]: [%d] rows for period [%s <-> %s] successfully moved",
				JobCount,
				move_count_row,
				time.Unix(move_timestamp_first, 0),
				time.Unix(move_timestamp_last, 0),
			)
		}

		log.Printf(conf.LOG_INFO+"Now sleep for [%d] seconds", conf.JOB_SLEEP_DURATION)
		time.Sleep(conf.JOB_SLEEP_DURATION * time.Second)
	}

	log.Printf(conf.LOG_WARN + "Too many failed jobs. No sense to continue")
	return nil
}
