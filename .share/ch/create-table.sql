/*  
  Example of table structure to mirror *CRTA-BOX* data in *ClickHouse*-database
*/

CREATE TABLE IF NOT EXISTS BOXes.box00 (
    
    -- Measurement localization
      -- in log
     `id`  UUID DEFAULT generateUUIDv4(1) COMMENT 'Record unique identifier formed in-situ'
    
      -- in time
    ,`timestamp` DateTime('Europe/Moscow')  NOT NULL  -- streamed tag *timestamp*
    
      -- in space
    ,`id_unit`         UInt64  NOT NULL  COMMENT '*CRTA-BOX* tag *unit_id* - boiler unit identification code provided by external register'
    ,`id_box`          UInt16  NOT NULL  COMMENT '*CRTA-BOX* tag *apk_id* - *CRTA-BOX* identifier provided by vendor'
    ,`id_spot`         UInt8   NOT NULL  COMMENT '*CRTA-BOX* tag *slave_id* - MODBUS standard register that serves as a spot identifier'
    ,`id_sensor`       UInt32  NOT NULL  COMMENT '*CRTA-BOX* tag *crta_id* (*device_id* in SrtaReg settings.ini) - identifier of signal sensor (of telltale) provided by external register'
    
      -- Sensor instance specifications
    ,`sensor_serial`   UInt32  NOT NULL  COMMENT '*CRTA-BOX* serial number provided by vendor'
    ,`sensor_revision` UInt8   NOT NULL  COMMENT '*CRTA-BOX* serial number provided by vendor'

    -- Acoustic measurements and diagnosis
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
        
        ,'fault_log'      = 127  -- extra-value formed by *crta-box* streaming process when original value is out of [0-8] range
      ) NOT NULL  COMMENT '*CRTA-BOX* tag *status_modbus* with value interpretation'
    ,`value_lf`        UInt16 NOT NULL COMMENT '*CRTA-BOX* rev.5 tag *UN_modbus*'
    ,`value_hf`        UInt16 NOT NULL COMMENT '*CRTA-BOX* rev.5 tag *UV_modbus*'
    ,`value_signal`    UInt16 NOT NULL COMMENT '*CRTA-BOX* rev.5 tag *Yi_modbus*'
    ,`value_ratio`      Int8  NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *Pi_modbus*'
    ,`preset_signal`   UInt8  NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *Yy_modbus*, absolute value'
    ,`preset_ratio`     Int8  NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *Py_modbus*'
    ,`threshold_halt`  UInt8  NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *UKOT_min_modbus*'
    ,`threshold_fault` UInt16 NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *UDAT_min_modbus*'
    
    -- *CRTA-BOX* introspection measurements
    ,`sensor_voltage` UInt8 NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *UMIK_modbus*'
    ,`relay_delay`    UInt8 NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *DELAY_modbus*'
    ,`filter_gain`    UInt8 NOT NULL COMMENT '*CRTA-BOX* rev 5 tag *AVG_filter_modbus*'

    -- ETL result status
    ,`StatusCode` UInt32  NOT NULL DEFAULT 0x00A20000 
      COMMENT '[OPC UA status code](http://www.opcfoundation.org/UA/schemas/StatusCode.csv) formed by *crta-box* streaming process'
) ENGINE = ReplacingMergeTree() -- No duplicates are allowed
PARTITION BY (id_unit, id_box, id_spot, toYYYYMM(`timestamp`))
ORDER BY (id_unit, id_box, id_spot, `timestamp`)
