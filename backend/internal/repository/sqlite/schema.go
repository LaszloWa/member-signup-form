package sqlite

const schemaSQL = `
CREATE TABLE IF NOT EXISTS submissions (
	submission_id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	phone_number TEXT NOT NULL,
	birth_date TEXT NOT NULL,
	member_type_id TEXT NOT NULL,
	club_id TEXT NOT NULL,
	form_id TEXT NOT NULL,
	created_at TEXT NOT NULL,
	email_normalized TEXT NOT NULL,
	phone_normalized TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_submissions_duplicate_scope
ON submissions (form_id, club_id, email_normalized, phone_normalized, birth_date);

CREATE INDEX IF NOT EXISTS idx_submissions_created_at
ON submissions (created_at);
`
