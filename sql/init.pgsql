-- Clean up previous instance

DROP TABLE IF EXISTS user_signup_request CASCADE;
DROP TABLE IF EXISTS password_reset_request CASCADE;
DROP TABLE IF EXISTS user_session CASCADE;
DROP TABLE IF EXISTS user_account CASCADE;
DROP TYPE IF EXISTS user_role_type;

DROP COLLATION IF EXISTS case_insensitive;

--------------------------------------------------
-- Create user management and session tables

CREATE TYPE user_role_type AS ENUM (
	'admin', -- can do anything
	'moderator', -- can delete and edit stuff
	'user', -- can create categories, posts, comments, and votes
	'inactive', -- can't do anything
	'banned' -- can't do anything
);

CREATE TABLE user_account (
	id SERIAL PRIMARY KEY,
	email VARCHAR(50) UNIQUE NOT NULL,
	user_role user_role_type NOT NULL DEFAULT 'user',
	handle VARCHAR(25) UNIQUE, -- optional handle
	display_name VARCHAR(50) NOT NULL, -- required
	auth_hash VARCHAR(60) NOT NULL,
	user_settings JSON,
	created_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX user_account_handle_idx ON user_account (handle);
CREATE UNIQUE INDEX user_account_email_idx ON user_account (email);

CREATE TABLE user_session (
	token VARCHAR(30) PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	expires TIMESTAMPTZ NOT NULL
);

CREATE TABLE user_signup_request (
	id SERIAL PRIMARY KEY,
	email VARCHAR(50) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	token VARCHAR(15) UNIQUE NOT NULL
);

CREATE TABLE password_reset_request (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	sent_to_address VARCHAR(50) NOT NULL,
	token VARCHAR(15) UNIQUE,
	created_at TIMESTAMPTZ
);

--------------------------------------------------
-- Create metadata objects

CREATE COLLATION case_insensitive (
	provider = icu, -- "International Components for Unicode"
	-- und stands for undefined (ICU root collation - language agnostic)
	-- colStrength=primary ignores case and accents
	-- colNumeric=yes sorts strings with numeric parts by numeric value
	-- colAlternate=shifted would recognize equality of equivalent punctuation sequences
	locale = 'und@colStrength=primary;colNumeric=yes',
	deterministic = false
);

--------------------------------------------------

CREATE TYPE space_type ENUM (
	'space', -- (nameless; contains titles and other spaces)
	'user', -- (user's personal space)
	'naked-tag', -- per-character-time-data text without a double newline
	'tag', -- plaintext without a double newline
	'title', -- plaintext with a newline
	'view', -- (tag intersection)
	'user-checkin', -- user checking themselves in personally to a space
	'space-checkin', -- user checking in a space to another space
	'json-attribute', -- URL and json path and refresh rate
	'picture',
	'audio',
	'video',
	'stream-of-consciousness',
	'posted-note'
);

CREATE TABLE space (
	id SERIAL PRIMARY KEY,
	space_type space_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE TABLE user_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE TABLE tag_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	tag_text TEXT COLLATE case_insensitive NOT NULL
);

CREATE TABLE naked_tag_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	tag_text TEXT NOT NULL,
	replay_data JSON NOT NULL
);

CREATE TABLE title_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	title_text TEXT NOT NULL,
	replay_data JSON
);

CREATE TABLE view_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	component_space_ids INTEGER[] NOT NULL
);

CREATE TYPE checkin_time (
	'past',
	'now',
	'future'
);

CREATE TABLE user_checkin_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	checkin_time checkin_time NOT NULL
);

CREATE TABLE space_checkin_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	checkin_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	checkin_time checkin_time NOT NULL
);

CREATE TABLE json_attribute_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	url TEXT NOT NULL,
	json_path TEXT NOT NULL,
	refresh_rate INTERVAL
);

CREATE TYPE stream_of_consciousness_mode (
	'naked',
	'clothed'
);

CREATE TABLE stream_of_consciousness_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	clothing stream_of_consciousness_mode NOT NULL,
	final_text TEXT NOT NULL
);

CREATE TABLE posted_note_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	note_text TEXT NOT NULL
);
