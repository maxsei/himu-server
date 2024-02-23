-- $ rg 'GLIBCXX_3.4.30' --binary -g '*libstdc\+\+*' -l /nix/store/
INSTALL sqlite3;
LOAD sqlite3;

CREATE VIEW v1 AS SELECT server_start_ts_ms,
                ts_ms,
                "value",
                Concat("sensor", ' ', ( CASE
                                          WHEN coord = 0 THEN 'x'
                                          WHEN coord = 1 THEN 'y'
                                          WHEN coord = 2 THEN 'z'
                                          ELSE ''
                                        END )) AS sensor
         FROM   Sqlite_scan("./data/sensor_data.sqlite3", "sensor_data")

-- PIVOT v1 ON "sensor" USING first("value");
