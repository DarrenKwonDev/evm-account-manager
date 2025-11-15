
-- sqlite3 app.db < db/migrations/001_init.up.sql

-- account table
CREATE TABLE IF NOT EXISTS accounts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	address TEXT NOT NULL UNIQUE,
	private_key TEXT NOT NULL UNIQUE,
	alias TEXT NOT NULL DEFAULT '',
	chain TEXT NOT NULL DEFAULT '',
	label TEXT,
	memo TEXT,
	total_value REAL DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP	
);

-- schema version
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO schema_migrations (version) VALUES ('001');


-- index
CREATE INDEX IF NOT EXISTS idx_accounts_label ON accounts(label);
CREATE INDEX IF NOT EXISTS idx_accounts_chain ON accounts(chain);
