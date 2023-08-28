package conf

const (
	// Aspects and presets
	GENESIS_TIMESTAMP  = 1685480400 // 2023-05-31 00:00:00 - date when first *CRTA-BOX* was installed
	JOB_MAX_FAIL_COUNT = 5
	JOB_SLEEP_DURATION = 250 // [seconds]

	// * tag *status_modbus* specifications (see ./crta-box/share/ch_create-tables.sql for details)
	TAG_PRESET_STATUS_MIN       = 0
	TAG_PRESET_STATUS_MAX       = 8
	TAG_PRESET_STATUS_INCORRECT = 127
)
