-- Clean up previous instance

DROP INDEX IF EXISTS space_time_idx;
DROP INDEX IF EXISTS space_type_time_idx;
DROP INDEX IF EXISTS space_user_throttle;
DROP TABLE IF EXISTS user_bookmark CASCADE;
DROP TABLE IF EXISTS user_space CASCADE;
DROP TABLE IF EXISTS user_space_config CASCADE;
DROP TABLE IF EXISTS json_attribute_space CASCADE;
DROP TABLE IF EXISTS stream_of_consciousness_space CASCADE;
DROP TABLE IF EXISTS text_space_revision CASCADE;
DROP TABLE IF EXISTS text_space CASCADE;
DROP TABLE IF EXISTS tag_space CASCADE;
DROP TABLE IF EXISTS link_space CASCADE;
DROP TABLE IF EXISTS branch_space CASCADE;
DROP TABLE IF EXISTS checkin CASCADE;
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
	'user', -- a user's personal space (always at root)
	'branch', -- has a label unique to its parent
	'text', -- plain text entered by a user, with optional replay data
	'tag', -- plain text (no newlines), special handling to give a space a set of active tags
	'link' -- user linking in a space to another space
	--'stream-of-consciousness', -- contains a stream of text checkins ("text-radio")
	--'json-attribute' -- URL and json path and refresh rate

	-- 'picture',
	-- 'audio',
	-- 'video',

	-- monetization -- fee is 1¢ per second
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
	created_by INTEGER NOT NULL REFERENCES user_account (id)
);

CREATE INDEX space_time_idx ON space (parent_id, created_at); -- for top queries
CREATE INDEX space_type_time_idx ON space (parent_id, space_type, created_at);
CREATE INDEX space_user_throttle ON space (created_by, created_at);

CREATE TABLE checkin (
	id SERIAL PRIMARY KEY,
	space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES user_account (id),
	created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE user_space ( -- a user's personal space (always at root)
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES user_account (id),
	UNIQUE (user_id)
);

CREATE TABLE user_space_config ( -- config for subspaces within a user space
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES user_account (id),
	order_number INTEGER NOT NULL
);

CREATE TABLE branch_space (
	space_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	parent_id INTEGER NOT NULL, -- 0 for root spaces
	label_text_id INTEGER NOT NULL REFERENCES unique_text (id),
	UNIQUE (parent_id, label_text_id)
);

CREATE TABLE link_space ( -- a link to another space somewhere else
	parent_id INTEGER NOT NULL REFERENCES space (id),
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	link_space_id INTEGER NOT NULL REFERENCES space (id),
	UNIQUE (parent_id, link_space_id)
);

CREATE TABLE tag_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	text_id INTEGER NOT NULL REFERENCES unique_text (id),
	UNIQUE (parent_id, text_id)
);

-- allow creating duplicate text spaces
CREATE TABLE text_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	title_id INTEGER REFERENCES unique_text (id), -- optional title for text space
	text_id INTEGER NOT NULL REFERENCES unique_text (id),
	recording TEXT, -- optional recording of text changes as JSON
	started_at TIMESTAMPTZ, -- null if no recording
	backdated_at TIMESTAMPTZ, -- null if no backdating
);

-- record history of changes
CREATE TABLE text_space_revision (
	space_id INTEGER NOT NULL REFERENCES text_space (space_id) ON DELETE CASCADE,
	revision_datetime TIMESTAMPTZ NOT NULL,
	title_id INTEGER REFERENCES unique_text (id), -- optional title for text space
	text_id INTEGER NOT NULL REFERENCES unique_text (id),
	PRIMARY KEY (space_id, revision_datetime)
);

CREATE TABLE stream_of_consciousness_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	stream_closed_at TIMESTAMPTZ -- null until closed (can be reopened)
);

CREATE TABLE json_attribute_space (
	space_id INTEGER PRIMARY KEY REFERENCES space (id) ON DELETE CASCADE,
	parent_id INTEGER NOT NULL REFERENCES space (id) ON DELETE CASCADE,
	url TEXT NOT NULL,
	json_path TEXT NOT NULL,
	refresh_rate INTERVAL,
	UNIQUE (parent_id, url, json_path, refresh_rate)
);

CREATE TABLE user_bookmark (
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
