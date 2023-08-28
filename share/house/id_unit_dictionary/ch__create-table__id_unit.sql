/* 
   Dictionary of power plant units
   Create source table for *id_unit_dictionary*
*/

CREATE TABLE IF NOT EXISTS CRTA.id_unit (
    id_unit      UInt64
   ,label_plant  FixedString(8)  NOT NULL    -- human understandable abbreviation of power plant
   ,number_unit  UInt8           NOT NULL    -- number of unit at power plant
) ENGINE = TinyLog

