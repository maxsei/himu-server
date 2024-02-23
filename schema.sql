create table if not exists sensor_data (
	server_start_ts_ms INTEGER NOT NULL,
	ts_ms  INTEGER NOT NULL,
	sensor TEXT NOT NULL,
	value  NUMBER NOT NULL,
	coord INTEGER NOT NULL -- offset in array from json output, usually this is xyz coord
);
