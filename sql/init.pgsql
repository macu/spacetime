-- Clean up previous instance

DROP TRIGGER IF EXISTS on_user_account_insert ON user_account;
DROP TRIGGER IF EXISTS tree_node_content_tsvector_update ON tree_node_content;

DROP FUNCTION IF EXISTS delete_user_content;
DROP FUNCTION IF EXISTS set_first_user_to_admin;
DROP FUNCTION IF EXISTS init_db_content;
DROP FUNCTION IF EXISTS update_text_search;
DROP FUNCTION IF EXISTS extract_text_from_json;

DROP INDEX IF EXISTS user_account_handle_idx;
DROP INDEX IF EXISTS user_account_email_idx;
DROP INDEX IF EXISTS tree_node_tag_vote_idx;
DROP INDEX IF EXISTS tree_node_tag_class_idx;
DROP INDEX IF EXISTS tree_node_tag_vote_user_idx;
DROP INDEX IF EXISTS tree_node_content_vote_idx;
DROP INDEX IF EXISTS tree_node_content_vote_user_idx;
DROP INDEX IF EXISTS tree_node_content_idx;
DROP INDEX IF EXISTS tree_node_content_tag_vote_idx;
DROP INDEX IF EXISTS tree_node_content_tag_class_idx;
DROP INDEX IF EXISTS tree_node_content_tag_vote_user_idx;
DROP INDEX IF EXISTS tree_node_vote_idx;
DROP INDEX IF EXISTS tree_node_vote_user_idx;
DROP INDEX IF EXISTS tree_node_parent_idx;
DROP INDEX IF EXISTS tree_node_class_idx;
DROP INDEX IF EXISTS tree_node_lang_code_node_idx;
DROP INDEX IF EXISTS tree_node_lang_code_idx;
DROP INDEX IF EXISTS tree_node_meta_idx;
DROP INDEX IF EXISTS tree_node_meta_node_idx;
DROP INDEX IF EXISTS tree_node_owner_ordering_idx;
DROP INDEX IF EXISTS tree_node_merge_vote_idx;
DROP INDEX IF EXISTS tree_node_merge_vote_user_idx;
DROP INDEX IF EXISTS tree_node_content_search_idx;
DROP INDEX IF EXISTS tree_node_content_search_composite_idx;

DROP TABLE IF EXISTS tree_node_tag_vote;
DROP TABLE IF EXISTS tree_node_content_tag_vote;
DROP TABLE IF EXISTS tree_node_content_vote;
DROP TABLE IF EXISTS tree_node_content;
DROP TABLE IF EXISTS tree_node_vote;
DROP TABLE IF EXISTS tree_node_lang_code;
DROP TABLE IF EXISTS tree_node_meta;
DROP TABLE IF EXISTS tree_node_owner_ordering;
DROP TABLE IF EXISTS tree_node_merge_vote;
DROP TABLE IF EXISTS tree_node;

DROP TABLE IF EXISTS user_signup_request;
DROP TABLE IF EXISTS password_reset_request;
DROP TABLE IF EXISTS user_session;
DROP TABLE IF EXISTS user_account;

DROP TYPE IF EXISTS user_role_type;
DROP TYPE IF EXISTS tree_node_class;
DROP TYPE IF EXISTS owner_type;
DROP TYPE IF EXISTS vote_type;
DROP TYPE IF EXISTS tree_node_content_type;
DROP TYPE IF EXISTS tree_node_body_type;

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
	'agree',
	'disagree'
);

CREATE TYPE owner_type AS ENUM (
	'admin',
	'public',
	'user'
);

-- create table for tree nodes
CREATE TYPE tree_node_class AS ENUM (
	'lang',
	'tag',
	'type',
	'field',
	'category',
	'post', -- posts are always user-owned
	'comment' -- comments are always user-owned
);
CREATE TABLE tree_node (
	id SERIAL PRIMARY KEY,
	parent_id INTEGER REFERENCES tree_node (id) ON DELETE CASCADE, -- allow null for root
	node_class tree_node_class NOT NULL,
	owned_by owner_type NOT NULL DEFAULT 'public',
	is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX tree_node_parent_idx ON tree_node (parent_id, is_deleted, owned_by);
CREATE INDEX tree_node_class_idx ON tree_node (parent_id, is_deleted, node_class);

-- create table for associating internal keys with specific system nodes
CREATE TABLE tree_node_meta (
	id SERIAL PRIMARY KEY,
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	internal_key VARCHAR(10) NOT NULL
);

CREATE UNIQUE INDEX tree_node_meta_node_idx ON tree_node_meta (node_id);
CREATE UNIQUE INDEX tree_node_meta_idx ON tree_node_meta (internal_key);

-- create table for ordering user-owned nodes
CREATE TABLE tree_node_owner_ordering (
	id SERIAL PRIMARY KEY,
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	order_number INTEGER NOT NULL
);

CREATE UNIQUE INDEX tree_node_owner_ordering_idx ON tree_node_owner_ordering (node_id);

-- create table for associating lang codes with lang tags
CREATE TABLE tree_node_lang_code (
	id SERIAL PRIMARY KEY,
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	lang_code VARCHAR(10) NOT NULL
);

CREATE INDEX tree_node_lang_code_node_idx ON tree_node_lang_code (node_id);
CREATE INDEX tree_node_lang_code_idx ON tree_node_lang_code (lang_code);

-- create table for voting on node placements
CREATE TABLE tree_node_vote (
	id SERIAL PRIMARY KEY,
	parent_id INTEGER REFERENCES tree_node (id) ON DELETE CASCADE, -- null for root
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE INDEX tree_node_vote_idx ON tree_node_vote (parent_id, node_id, vote);
CREATE UNIQUE INDEX tree_node_vote_user_idx ON tree_node_vote (created_by, parent_id, node_id);

-- create table for node title/body content
CREATE TYPE tree_node_content_type AS ENUM (
	'title', -- categories, tags, types, fields, langs, posts
	'description', -- categories
	'comment-body', -- comments
	'post-json' -- block on posts
);
CREATE TABLE tree_node_content (
	id SERIAL PRIMARY KEY,
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	content_type tree_node_content_type NOT NULL,
	text_content TEXT NOT NULL COLLATE case_insensitive,
	text_search TSVECTOR, -- for searching text_content
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX tree_node_content_idx ON tree_node_content (node_id, content_type, text_content);
CREATE INDEX tree_node_content_search_idx ON tree_node_content USING GIN (text_search);
CREATE INDEX tree_node_content_search_composite_idx ON tree_node_content (content_type, text_search);

-- write functions to update tsvector on insert/update
CREATE OR REPLACE FUNCTION extract_text_from_json(json_data JSONB) RETURNS TEXT AS $$
DECLARE
	extracted_text TEXT;
BEGIN
	-- Extract text from the JSON body.
	SELECT string_agg(value, ' ') INTO extracted_text
	FROM jsonb_each_text(json_data)
	WHERE key IN ('text');

	RETURN extracted_text;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_text_search() RETURNS TRIGGER AS $$
DECLARE
	text_to_index TEXT;
BEGIN
	IF NEW.content_type = 'post-json' THEN
		text_to_index := extract_text_from_json(NEW.text_content::JSONB);
	ELSE
		text_to_index := NEW.text_content;
	END IF;

	IF text_to_index IS NULL OR text_to_index = '' THEN
		NEW.text_search := NULL;
	ELSE
		NEW.text_search := to_tsvector('pg_catalog.simple', text_to_index);
	END IF;

	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tree_node_content_tsvector_update BEFORE INSERT OR UPDATE
ON tree_node_content FOR EACH ROW EXECUTE FUNCTION update_text_search();

-- create table for voting on node title/body content
CREATE TABLE tree_node_content_vote (
	id SERIAL PRIMARY KEY,
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	content_id INTEGER NOT NULL REFERENCES tree_node_content (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE INDEX tree_node_content_vote_idx ON tree_node_content_vote (node_id, content_id, vote);
CREATE UNIQUE INDEX tree_node_content_vote_user_idx ON tree_node_content_vote (created_by, node_id, content_id);

-- create table for tagging notes
CREATE TABLE tree_node_tag_vote (
	id SERIAL PRIMARY KEY,
	node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	tag_node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	tag_class tree_node_class NOT NULL,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE INDEX tree_node_tag_vote_idx ON tree_node_tag_vote (node_id, tag_node_id, vote);
CREATE INDEX tree_node_tag_class_idx ON tree_node_tag_vote (node_id, tag_class, tag_node_id, vote);
CREATE UNIQUE INDEX tree_node_tag_vote_user_idx ON tree_node_tag_vote (created_by, node_id, tag_node_id);

-- create table for tagging node content
CREATE TABLE tree_node_content_tag_vote (
	id SERIAL PRIMARY KEY,
	content_id INTEGER NOT NULL REFERENCES tree_node_content (id) ON DELETE CASCADE,
	tag_node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	tag_class tree_node_class NOT NULL,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE INDEX tree_node_content_tag_vote_idx ON tree_node_content_tag_vote (content_id, tag_node_id, vote);
CREATE INDEX tree_node_content_tag_class_idx ON tree_node_content_tag_vote (content_id, tag_class, tag_node_id, vote);
CREATE INDEX tree_node_content_tag_vote_user_idx ON tree_node_content_tag_vote (created_by, content_id, tag_node_id);

-- create table for soft merging nodes
CREATE TABLE tree_node_merge_vote (
	id SERIAL PRIMARY KEY,
	target_node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	source_node_id INTEGER NOT NULL REFERENCES tree_node (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE
);

CREATE INDEX tree_node_merge_vote_idx ON tree_node_merge_vote (target_node_id, source_node_id, vote);
CREATE UNIQUE INDEX tree_node_merge_vote_user_idx ON tree_node_merge_vote (created_by, target_node_id, source_node_id);

--------------------------------------------------
-- create triggers

--------------------------------------------------
-- create functions

-- delete everything assocaited with the user, except categories and tags
CREATE FUNCTION delete_user_content (
	user_id_param INTEGER
)
RETURNS VOID
LANGUAGE PLPGSQL
AS $$
BEGIN
	DELETE FROM user_session WHERE user_id = user_id_param;
	DELETE FROM user_account WHERE id = user_id;
	DELETE FROM password_reset_request WHERE user_id = user_id_param;
	-- TODO delete content across tree
END;
$$;

--------------------------------------------------
-- create db init function

CREATE FUNCTION init_db_content (
	user_id INTEGER
)
RETURNS VOID
LANGUAGE PLPGSQL
AS $$
DECLARE
	root_node_id INTEGER;
	langs_node_id INTEGER;
	lang_en_node_id INTEGER;
	tags_node_id INTEGER;
	types_node_id INTEGER;
BEGIN
	-- exit if any tree node exists
	PERFORM 1 FROM tree_node LIMIT 1;
	IF FOUND THEN
		RETURN;
	END IF;

	-- create root node
	INSERT INTO tree_node (node_class, owned_by, created_at, created_by)
	VALUES ('category', 'admin', NOW(), user_id)
	RETURNING id INTO root_node_id;

	INSERT INTO tree_node_content (node_id, content_type, text_content, created_at, created_by)
	VALUES (root_node_id, 'title', 'TreeTime', NOW(), user_id);

	INSERT INTO tree_node_meta (node_id, internal_key)
	VALUES (root_node_id, 'treetime');

	-- create languages category
	INSERT INTO tree_node (parent_id, node_class, owned_by, created_at, created_by)
	VALUES (root_node_id, 'category', 'admin', NOW(), user_id)
	RETURNING id INTO langs_node_id;

	INSERT INTO tree_node_content (node_id, content_type, text_content, created_at, created_by)
	VALUES (langs_node_id, 'title', 'Languages', NOW(), user_id);

	INSERT INTO tree_node_meta (node_id, internal_key)
	VALUES (langs_node_id, 'langs');

	-- create English lang
	INSERT INTO tree_node (parent_id, node_class, owned_by, created_at, created_by)
	VALUES (langs_node_id, 'lang', 'admin', NOW(), user_id)
	RETURNING id INTO lang_en_node_id;

	INSERT INTO tree_node_content (node_id, content_type, text_content, created_at, created_by)
	VALUES (lang_en_node_id, 'title', 'English', NOW(), user_id);

	INSERT INTO tree_node_lang_code (node_id, lang_code)
	VALUES (lang_en_node_id, 'en');

	-- create tags category
	INSERT INTO tree_node (parent_id, node_class, owned_by, created_at, created_by)
	VALUES (root_node_id, 'category', 'admin', NOW(), user_id)
	RETURNING id INTO tags_node_id;

	INSERT INTO tree_node_content (node_id, content_type, text_content, created_at, created_by)
	VALUES (tags_node_id, 'title', 'Tags', NOW(), user_id);

	INSERT INTO tree_node_meta (node_id, internal_key)
	VALUES (tags_node_id, 'tags');

	-- create types category
	INSERT INTO tree_node (parent_id, node_class, owned_by, created_at, created_by)
	VALUES (root_node_id, 'category', 'admin', NOW(), user_id)
	RETURNING id INTO types_node_id;

	INSERT INTO tree_node_content (node_id, content_type, text_content, created_at, created_by)
	VALUES (types_node_id, 'title', 'Types', NOW(), user_id);

	INSERT INTO tree_node_meta (node_id, internal_key)
	VALUES (types_node_id, 'types');
END;
$$;

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

	PERFORM init_db_content(NEW.id);

	-- delete trigger and function
	DROP TRIGGER on_user_account_insert ON user_account;
	DROP FUNCTION set_first_user_to_admin;
	DROP FUNCTION init_db_content;

	RETURN NULL;
END;
$$;

CREATE TRIGGER on_user_account_insert AFTER INSERT ON user_account
	FOR EACH ROW EXECUTE PROCEDURE set_first_user_to_admin();
