package conf

const (
	// SQL-queries and templates

	// - *ClickHouse* queries
	SQL_SELECT_HOUSE_MAX_TIMESTAMP = "SELECT toUInt64(MAX(timestamp)) FROM log_box3"
	SQL_INSERT_HOUSE_DATA          = "INSERT INTO CRTA.log_box3"

	// - *CRTA-BOX* (PostgreSQL) queries
	// -- Tag renaming in accordance with specification in *share/ch__create-table__log_box3.sql*
	SQL_SELECT_BOX_DATA = `
	SELECT 
		 timestamp
		,unit_id			AS id_unit
		,apk_id				AS id_box
		,slave_id			AS id_spot
		,crta_id			AS id_sensor
		,status_modbus		AS value_status
		,p100_modbus		AS value_ratio
		,u100_modbus		AS value_grade
		,ustp_modbus		AS preset_ratio
		,ustu_modbus		AS preset_grade
		,sko_modbus			AS value_noise
		,ustskomin_modbus	AS preset_noise
		,r_modbus			AS value_gain
		,verh_modbus		AS value_high
		,niz_modbus			AS value_low
	FROM cache_log_srta33_view 
		WHERE 
				timestamp > %d 
			AND timestamp < (
				SELECT MAX(timestamp) - 1 FROM log_srta33 WHERE timestamp >= %d
			)
	ORDER BY timestamp
	`
)
