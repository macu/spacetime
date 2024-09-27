-- Clean up previous instance

DROP TRIGGER IF EXISTS on_category_delete ON category;
DROP TRIGGER IF EXISTS on_post_delete ON post;
DROP TRIGGER IF EXISTS on_comment_delete ON comment;
DROP TRIGGER IF EXISTS on_tag_delete ON tag;

DROP FUNCTION IF EXISTS delete_user_content;

DROP INDEX IF EXISTS tag_vote_idx;
DROP INDEX IF EXISTS category_merge_vote_idx;
DROP INDEX IF EXISTS category_post_vote_idx;
DROP INDEX IF EXISTS category_title_vote_idx;
DROP INDEX IF EXISTS category_location_vote_idx;
DROP INDEX IF EXISTS category_title_idx;

DROP TABLE IF EXISTS category_merge_vote;
DROP TABLE IF EXISTS tag_vote;
DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS category_post_vote;
DROP TABLE IF EXISTS category_title_vote;
DROP TABLE IF EXISTS category_location_vote;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS category_title;
DROP TABLE IF EXISTS category;

DROP TABLE IF EXISTS user_signup_request;
DROP TABLE IF EXISTS password_reset_request;
DROP TABLE IF EXISTS user_session;
DROP TABLE IF EXISTS user_account;

DROP COLLATION IF EXISTS case_insensitive;

--------------------------------------------------
-- Create user management and session tables

CREATE TYPE user_role_type AS ENUM (
	'admin', -- can do anything
	'moderator', -- can delete and edit stuff
	'contributor' -- can create categories, posts, comments, and votes
);

CREATE TABLE user_account (
	id SERIAL PRIMARY KEY,
	email VARCHAR(50) UNIQUE NOT NULL,
	user_role user_role_type NOT NULL DEFAULT 'contributor',
	handle VARCHAR(25) UNIQUE COLLATE case_insensitive, -- optional handle
	display_name VARCHAR(50) NOT NULL, -- required
	auth_hash VARCHAR(60) NOT NULL,
	user_settings JSON,
	created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE user_session (
	token VARCHAR(30) PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE,
	expires TIMESTAMPTZ NOT NULL
);

CREATE TABLE user_signup_request (
	id SERIAL PRIMARY KEY,
	username VARCHAR(25) UNIQUE NOT NULL COLLATE case_insensitive,
	email VARCHAR(50) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	token VARCHAR(15) UNIQUE,
	user_id INTEGER -- set after email verification
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

CREATE TABLE category (
	id SERIAL PRIMARY KEY,
);

CREATE TABLE category_title (
	id SERIAL PRIMARY KEY,
	category_id INTEGER NOT NULL REFERENCES category (id) ON DELETE CASCADE,
	local_parent_id INTEGER REFERENCES category (id) ON DELETE CASCADE,
	lang VARCHAR(5) NOT NULL,
	title VARCHAR(128) NOT NULL COLLATE case_insensitive,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX category_title_idx ON category_title (category_id, local_parent_id, lang, title);

CREATE TABLE post (
	id SERIAL PRIMARY KEY,
	lang VARCHAR(5) NOT NULL,
	title VARCHAR(128), -- optional
	body TEXT NOT NULL, -- required
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL,
	updated_at TIMESTAMPTZ,
	updated_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE TABLE tag (
	id SERIAL PRIMARY KEY,
	title VARCHAR(64) UNIQUE NOT NULL COLLATE case_insensitive,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE TYPE vote_type AS ENUM (
	'agree',
	'disagree'
);

CREATE TABLE category_location_vote (
	id SERIAL PRIMARY KEY,
	parent_id INTEGER REFERENCES category (id) ON DELETE CASCADE,
	child_id INTEGER REFERENCES category (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX category_location_vote_idx ON category_location_vote (parent_id, child_id, vote);

CREATE TABLE category_title_vote (
	id SERIAL PRIMARY KEY,
	category_id INTEGER REFERENCES category (id) ON DELETE CASCADE,
	category_title_id INTEGER REFERENCES category_title (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX category_title_vote_idx ON category_title_vote (category_id, category_title_id, vote);

CREATE TABLE category_post_vote (
	id SERIAL PRIMARY KEY,
	category_id INTEGER REFERENCES category (id) ON DELETE CASCADE,
	post_id INTEGER REFERENCES post (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX category_post_vote_idx ON category_post_vote (category_id, post_id, vote);

CREATE TYPE comment_target_type AS ENUM (
	'category',
	'post'
);

CREATE TABLE comment (
	id SERIAL PRIMARY KEY,
	target_type comment_target_type NOT NULL,
	target_id INTEGER NOT NULL,
	body TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL,
	updated_at TIMESTAMPTZ,
	updated_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE TYPE tag_target_type AS ENUM (
	'category',
	'category_title',
	'post',
	'comment'
);

CREATE TABLE tag_vote (
	id SERIAL PRIMARY KEY,
	tag_id INTEGER REFERENCES tag (id) ON DELETE CASCADE,
	target_type tag_target_type NOT NULL,
	target_id INTEGER,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX tag_vote_idx ON tag_vote (target_type, target_id, tag_id, vote);

CREATE TABLE category_merge_vote (
	id SERIAL PRIMARY KEY,
	source_id INTEGER REFERENCES category (id) ON DELETE CASCADE,
	target_id INTEGER REFERENCES category (id) ON DELETE CASCADE,
	vote vote_type NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by INTEGER REFERENCES user_account (id) ON DELETE SET NULL
);

CREATE INDEX category_merge_vote_idx ON category_merge_vote (target_id, source_id, vote);

--------------------------------------------------
-- create triggers

CREATE FUNCTION delete_category ()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
BEGIN
	DELETE FROM tag_vote WHERE target_type = 'category' AND target_id = OLD.id;
	DELETE FROM comment WHERE target_type = 'category' AND target_id = OLD.id;
	RETURN NULL;
END $$;

CREATE TRIGGER on_category_delete AFTER DELETE ON category
	FOR EACH ROW EXECUTE PROCEDURE delete_category();

CREATE FUNCTION delete_post ()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
BEGIN
	DELETE FROM tag_vote WHERE target_type = 'post' AND target_id = OLD.id;
	DELETE FROM comment WHERE target_type = 'post' AND target_id = OLD.id;
	RETURN NULL;
END $$;

CREATE TRIGGER on_post_delete AFTER DELETE ON post
	FOR EACH ROW EXECUTE PROCEDURE delete_post();

CREATE FUNCTION delete_comment ()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
BEGIN
	DELETE FROM tag_vote WHERE target_type = 'comment' AND target_id = OLD.id;
	RETURN NULL;
END $$;

CREATE TRIGGER on_comment_delete AFTER DELETE ON comment
	FOR EACH ROW EXECUTE PROCEDURE delete_comment();

CREATE FUNCTION delete_tag ()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
BEGIN
	DELETE FROM tag_vote WHERE tag_id = OLD.id;
	RETURN NULL;
END $$;

CREATE TRIGGER on_tag_delete AFTER DELETE ON tag
	FOR EACH ROW EXECUTE PROCEDURE delete_tag();

--------------------------------------------------
-- create functions

-- delete everything assocaited with the user, except categories and tags
CREATE FUNCTION delete_user_content (
	user_id INTEGER
) RETURNS VOID AS $$
BEGIN
	DELETE FROM user_session WHERE user_id = user_id;
	DELETE FROM user_account WHERE id = user_id;
	DELETE FROM user_signup_request WHERE user_id = user_id;
	DELETE FROM password_reset_request WHERE user_id = user_id;
	DELETE FROM category_location_vote WHERE created_by = user_id;
	DELETE FROM category_title_vote WHERE created_by = user_id;
	DELETE FROM category_post_vote WHERE created_by = user_id;
	DELETE FROM post WHERE created_by = user_id;
	DELETE FROM comment WHERE created_by = user_id;
	DELETE FROM tag_vote WHERE created_by = user_id;
	DELETE FROM category_merge_vote WHERE created_by = user_id;
END;
