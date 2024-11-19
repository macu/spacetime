-- Clean up previous instance

DROP TRIGGER IF EXISTS on_user_account_insert ON user_account;

DROP FUNCTION IF EXISTS set_first_user_to_admin;

DROP TABLE IF EXISTS tag_vote_status CASCADE;
DROP TABLE IF EXISTS tag CASCADE;
DROP TABLE IF EXISTS tag_vote CASCADE;
DROP TABLE IF EXISTS tree_node_graft_vote_status CASCADE;
DROP TABLE IF EXISTS tree_node_graft_vote CASCADE;
DROP TABLE IF EXISTS tree_node_text_version CASCADE;
DROP TABLE IF EXISTS tree_node CASCADE;

DROP TYPE IF EXISTS tree_node_class;
DROP TYPE IF EXISTS vote_type;

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
-- Create content tables

CREATE TYPE vote_type AS ENUM (
	'up',
	'down'
);

CREATE TYPE tree_node_class AS ENUM (
	'plaintext',
	'markdown',
	'html',
	'image'
);

CREATE TABLE tree_node (
	id SERIAL PRIMARY KEY,
	parent_id INTEGER REFERENCES tree_node (id) ON DELETE CASCADE, -- null for root
	node_class tree_node_class NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX tree_node_parent_idx ON tree_node (parent_id);

CREATE TABLE tree_node_graft_vote (
	parent_id INTEGER REFERENCES tree_node (id) ON DELETE CASCADE, -- null for root
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	PRIMARY KEY (parent_id, node_id, created_by)
);

CREATE INDEX tree_node_graft_vote_idx ON tree_node_graft_vote (parent_id, node_id, vote);

CREATE TABLE tree_node_graft_vote_status (
	parent_id INTEGER REFERENCES tree_node (id) ON DELETE CASCADE, -- null for root
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	up_votes INTEGER NOT NULL DEFAULT 0,
	down_votes INTEGER NOT NULL DEFAULT 0,
	sum INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (parent_id, node_id)
);

CREATE TABLE tree_node_text_version (
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	text_content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	PRIMARY KEY (node_id, created_at)
);

CREATE TABLE tag (
	id SERIAL PRIMARY KEY,
	parent_id INTEGER REFERENCES tag (id) ON DELETE CASCADE, -- null for root
	tag_name VARCHAR(50) COLLATE case_insensitive NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id)
);

CREATE TABLE tag_vote (
	tag_id INTEGER NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	PRIMARY KEY (tag_id, created_by)
);

CREATE INDEX tree_node_tag_vote_idx ON tag_vote (tag_id, vote);

CREATE TABLE tag_vote_status (
	tag_id INTEGER NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
	up_votes INTEGER NOT NULL DEFAULT 0,
	down_votes INTEGER NOT NULL DEFAULT 0,
	sum INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (tag_id)
);

--------------------------------------------------
-- set first user to admin

CREATE FUNCTION set_first_user_to_admin ()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
BEGIN
	PERFORM 1 FROM user_account WHERE user_role = 'admin' LIMIT 1;
	IF FOUND THEN
		RETURN NULL;
	END IF;

	NEW.user_role := 'admin';

	-- delete trigger and function
	DROP TRIGGER on_user_account_insert ON user_account;
	DROP FUNCTION set_first_user_to_admin;

	RETURN NULL;
END;
$$;

CREATE TRIGGER on_user_account_insert AFTER INSERT ON user_account
	FOR EACH ROW EXECUTE PROCEDURE set_first_user_to_admin();
