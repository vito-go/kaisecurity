CREATE TABLE IF NOT EXISTS scan_results_info (
  scan_id TEXT PRIMARY KEY NOT NULL,
  timestamp TEXT,
  scan_status TEXT,
  resource_type TEXT,
  resource_name TEXT,
  summary TEXT,
  scan_metadata TEXT,
  create_time TEXT DEFAULT CURRENT_TIMESTAMP,
  update_time TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_scan_id ON scan_results_info ("scan_id");

CREATE TABLE IF NOT EXISTS vulnerability (
     scan_id TEXT NOT NULL,
     id TEXT NOT NULL,
     severity TEXT,
     cvss REAL,
     status TEXT,
     package_name TEXT,
     current_version TEXT,
     fixed_version TEXT,
     description TEXT,
     published_date TEXT,
     link TEXT,
     risk_factors TEXT,
     create_time TEXT DEFAULT CURRENT_TIMESTAMP,
     update_time TEXT DEFAULT CURRENT_TIMESTAMP
    );
CREATE UNIQUE INDEX IF NOT EXISTS idx_scan_id_id ON scan_results_info ("scan_id","id");
CREATE INDEX IF NOT EXISTS idx_severity_update_time ON vulnerability(severity, update_time);