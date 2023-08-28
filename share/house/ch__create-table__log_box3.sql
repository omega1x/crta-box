/* Create table for storing *CRTA-BOX* data in *ClickHouse* */

CREATE TABLE IF NOT EXISTS CRTA.log_box3 (
    
    -- Measurement localization in log
     `id`  UUID DEFAULT generateUUIDv4(1) COMMENT 'Record unique identifier formed in-situ'
    
    -- measurement localization in time
    ,`timestamp` DateTime('Europe/Moscow')  NOT NULL  -- streamed tag *timestamp*
    
    -- Measurement localization in space
      --- ordering  with
    ,`id_unit`   UInt64  NOT NULL  COMMENT '*CRTA-BOX* tag *unit_id*  - boiler identification code'
    ,`id_box`    UInt16  NOT NULL  COMMENT '*CRTA-BOX* tag *apk_id*   - *CRTA-BOX* identifier (odd values are associated with boiler unit A, even values - with B)'
    ,`id_spot`   UInt8   NOT NULL  COMMENT '*CRTA-BOX* tag *slave_id* - MODBUS standard register that serves as a spot identifier'
      ---
    ,`id_sensor` UInt32  NOT NULL  COMMENT '*CRTA-BOX* tag *crta_id*  - identifier of signal sensor (of telltale)'

    -- Measurements
    ,`value_status` Enum8(    
         'boiler_up'      =   0
        ,'fault_shortage' =   1
        ,'fault_unlink'   =   2
        ,'fault_sensor'   =   3
        ,'boiler_down'    =   4
        ,'rupture_param'  =   5
        ,'rupture_level'  =   6
        ,'rupture_multi'  =   7
        ,'device_care'    =   8  -- undocumented option 
        
        ,'fault_log'      = 127  -- extra-value formed by *crta-box* streaming process when original value is out of [0-7] range
      ) NOT NULL  COMMENT '*CRTA-BOX* tag *status_modbus* with value interpretation'
    ,`value_ratio`  Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *p100_modbus*'
    ,`value_grade`  Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *u100_modbus*'
    ,`preset_ratio` Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *ustp_modbus*'
    ,`preset_grade` Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *ustu_modbus*' 
    ,`value_noise`  Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *sko_modbus*'
    ,`preset_noise` Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *ustskomin_modbus*'
    ,`value_gain`   Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *r_modbus*'
    ,`value_high`   Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *verh_modbus*'
    ,`value_low`    Int16  NOT NULL  COMMENT '*CRTA-BOX* tag *niz_modbus*'

    -- Quality
    ,`StatusCode` UInt32  NOT NULL DEFAULT 0x00A20000 
      COMMENT '[OPC UA status code](http://www.opcfoundation.org/UA/schemas/StatusCode.csv) formed by *crta-box* streaming process'
) ENGINE = ReplacingMergeTree() -- No duplicates are allowed
PARTITION BY (id_unit, id_box, id_spot, toYYYYMM(`timestamp`))
ORDER BY (id_unit, id_box, id_spot, `timestamp`)


