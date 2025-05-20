/*
Package toolkit implements utility routines for data connectivity
between source and target databases
*/

package toolkit

const (
	BOX_TABLE_ORIGIN           = "box_view"
	BOX_DATA_EXPIRY_PERIOD     = 31536000
	SQL_TEMPLATE_MAX_TIMESTAMP = "SELECT MAX(timestamp) FROM %s"

	DEFAULT_HOUSE_TIME_LAG = -1200      // [seconds].
	GENESIS_TIMESTAMP      = 1685480400 // 2023-05-31 00:00:00 - date when first *CRTA-BOX* was installed
	BOX_DATABASE_SYSTEM    = "postgres"

	OPC_CODE_GoodEntryInserted = 0x00A20000 //
	OPC_CODE_BadOutOfRange     = 0x803C0000

	PRESET_MIN_VALUE_STATUS       = 0
	PRESET_MAX_VALUE_STATUS       = 8
	PRESET_INCORRECT_VALUE_STATUS = 127 // `fault_log` label for Enum8 value_status in *ClickHouse* (see [.share/ch/create-table.sql])
)
