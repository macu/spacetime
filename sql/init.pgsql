-- Clean up previous instance

DROP INDEX IF EXISTS space_time_idx;
DROP INDEX IF EXISTS space_type_time_idx;
DROP INDEX IF EXISTS space_user_throttle;
DROP INDEX IF EXISTS idx_user_space_user_id;
DROP TABLE IF EXISTS user_space_bookmark CASCADE;
DROP TABLE IF EXISTS user_space CASCADE;
DROP TABLE IF EXISTS space_checkin_activity CASCADE;
DROP TABLE IF EXISTS json_attribute_space CASCADE;
DROP TABLE IF EXISTS naked_text_space CASCADE;
DROP TABLE IF EXISTS stream_of_consciousness_space CASCADE;
DROP TABLE IF EXISTS text_space CASCADE;
DROP TABLE IF EXISTS tag_space CASCADE;
DROP TABLE IF EXISTS title_space CASCADE;
DROP TABLE IF EXISTS checkin_space CASCADE;
DROP TABLE IF EXISTS space CASCADE;
DROP TABLE IF EXISTS unique_text CASCADE;
DROP TYPE IF EXISTS space_type;

DROP INDEX IF EXISTS user_account_handle_idx;
DROP INDEX IF EXISTS user_account_email_idx;
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

CREATE TABLE unique_text (
	id SERIAL PRIMARY KEY,
	text_value TEXT COLLATE case_insensitive UNIQUE NOT NULL
);

CREATE TYPE space_type AS ENUM (
	'user', -- user's personal space
	'space', -- (nameless; contains titles and other spaces)
	'check-in', -- user checking in directly on a space
	'space-link', -- user linking in a space to another space
	'title', -- plain text (no newlines), special handling to give a space an active title
	'tag', -- plain text (no newlines), special handling to give a space a set of active tags
	'text', -- plain text entered by a user
	'naked-text', -- text with realtime replay data
	'stream-of-consciousness', -- contains a stream of text checkins ("text-radio")
	'json-attribute' -- URL and json path and refresh rate

	-- 'picture',
	-- 'audio',
	-- 'video',

	-- monetization -- fee is 1¢. per second
	-- 'rental-space',
	-- 'rental-payment',
	-- 'rental-donor',
	-- 'rental-payee',
	-- 'rental-payout'
);

CREATE TABLE space ( -- a domain that contains subspaces
	id SERIAL PRIMARY KEY,
	parent_id INTEGER REFERENCES space (id) ON DELETE CASCADE,
	space_type space_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE INDEX space_time_idx ON space (parent_id, created_at); -- for top queries
CREATE INDEX space_type_time_idx ON space (parent_id, space_type, created_at);
CREATE INDEX space_user_throttle ON space (created_by, created_at);

CREATE TABLE user_space ( -- a user's personal space
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_space_user_id ON user_space (user_id);

CREATE TABLE link_space ( -- a link to another space somewhere else
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	link_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	UNIQUE (parent_space_id, link_space_id)
);

CREATE TABLE title_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	unique_text_id INTEGER NOT NULL REFERENCES unique_text (id) ON DELETE CASCADE,
	UNIQUE (parent_space_id, unique_text_id)
);

CREATE TABLE tag_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	unique_text_id INTEGER NOT NULL REFERENCES unique_text (id) ON DELETE CASCADE,
	UNIQUE (parent_space_id, unique_text_id)
);

CREATE TABLE text_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	unique_text_id INTEGER NOT NULL REFERENCES unique_text (id) ON DELETE CASCADE,
	UNIQUE (parent_space_id, unique_text_id)
);

CREATE TABLE naked_text_space (
	-- allow duplicates of final text (replay data will probably always be unique)
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	final_unique_text_id INTEGER NOT NULL REFERENCES unique_text (id) ON DELETE CASCADE,
	replay_data TEXT NOT NULL,
	typing_started_at TIMESTAMPTZ
);

CREATE TABLE stream_of_consciousness_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	stream_closed_at TIMESTAMPTZ -- null until closed
);

CREATE TABLE json_attribute_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	url TEXT NOT NULL,
	json_path TEXT NOT NULL,
	refresh_rate INTERVAL,
	UNIQUE (parent_space_id, url, json_path, refresh_rate)
);

CREATE TABLE user_space_bookmark (
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL,
	PRIMARY KEY (user_id, space_id)
);

-- CREATE TYPE rental_space_payout_type (
-- 	'creators',
-- 	'public',
-- 	'platform', -- basically a donation to my company
-- 	'none'
-- );

-- CREATE TABLE rental_space (
-- 	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
-- 	creator_ids INTEGER[] NOT NULL,
-- 	payout_type rental_space_payout_type NOT NULL,
-- 	private BOOLEAN NOT NULL DEFAULT FALSE, -- platform payouts cannot be private
-- 	approved BOOLEAN NOT NULL, -- private spaces must be approved before publishing
-- 	release_payment BOOLEAN NOT NULL DEFAULT FALSE
-- );

-- CREATE TABLE rental_space_payee (
-- 	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
-- 	payee_id INTEGER NOT NULL
-- );
