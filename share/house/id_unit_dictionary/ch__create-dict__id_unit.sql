/* 
   Dictionary of power plant units
   Create dictionary *id_unit_dictionary* from source table *id_unit*

   NOTE:
     For tag semantics see CREATE-statement of source table
*/

CREATE OR REPLACE DICTIONARY CRTA.id_unit_dictionary (
    id_unit      UInt64
   ,label_plant  String
   ,number_unit  UInt8
)
PRIMARY KEY id_unit
SOURCE(CLICKHOUSE(TABLE 'id_unit'))
LAYOUT(FLAT())
LIFETIME(300)
