DROP TABLE IF EXISTS "forums";

DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "posts";

DROP TABLE IF EXISTS "threads";

DROP TABLE IF EXISTS "votes";

DROP TABLE IF EXISTS "members";

CREATE TABLE "forums" (
  "fID"     SERIAL PRIMARY KEY,
  "posts"   INT8 DEFAULT '0',
  "slug"    CITEXT,
  "threads" INT4 DEFAULT '0',
  "title"   TEXT,
  "author"  TEXT
);

DROP INDEX IF EXISTS index_on_forums_slug;

CREATE UNIQUE INDEX  index_on_forums_slug
  ON forums (slug);


CREATE TABLE "users" (
  "uID"      SERIAL PRIMARY KEY,
  "email"    CITEXT,
  "nickname" CITEXT,
  "fullname" TEXT,
  "about"    TEXT
);

DROP INDEX IF EXISTS index_on_users_email;

CREATE UNIQUE INDEX index_on_users_email
  ON users (email);

DROP INDEX IF EXISTS index_on_users_nickname;

CREATE UNIQUE INDEX index_on_users_nickname
  ON users (nickname);

DROP INDEX IF EXISTS index_on_users_nickname_and_email;

CREATE UNIQUE INDEX index_on_users_nickname_and_email
  ON users (nickname, email);


CREATE TABLE "posts" (
  "pID"      SERIAL PRIMARY KEY,
  "author"   TEXT,
  "created"  TIMESTAMP WITH TIME ZONE DEFAULT now(),
  "forum"    TEXT,
  "isEdited" BOOLEAN DEFAULT 'false',
  "message"  TEXT,
  "parent"   INT8 DEFAULT '0',
  "thread"   BIGINT,
  "path"     INT8 []
);

DROP INDEX IF EXISTS index_on_posts_thread;

CREATE INDEX index_on_posts_thread
  ON posts (thread);

DROP INDEX IF EXISTS index_on_posts_parent;

CREATE INDEX index_on_posts_parent
  ON posts (parent);

DROP INDEX IF EXISTS index_on_posts_path;

CREATE INDEX index_on_posts_path
  ON posts USING GIN (path);

-- 

-- DROP INDEX IF EXISTS index_on_posts_id_and_thread;

-- CREATE INDEX index_on_posts_id_and_thread
--   ON posts ("pID", thread);

-- DROP INDEX IF EXISTS index_on_posts_path_and_thread;

-- CREATE INDEX index_on_posts_path_and_thread
--   ON posts (thread, path);

DROP INDEX IF EXISTS index_on_posts_path_and_thread_and_parent;

CREATE INDEX index_on_posts_path_and_thread_and_parent
  ON posts (path, thread, parent);

-- DROP INDEX IF EXISTS index_on_posts_id_and_path;

-- CREATE INDEX index_on_posts_id_and_path
--   ON posts ("pID", path);


CREATE TABLE "threads" (
  "tID"     SERIAL PRIMARY KEY,
  "author"  TEXT,
  "created" TIMESTAMP WITH TIME ZONE DEFAULT now(),
  "forum"   CITEXT,
  "message" TEXT,
  "slug"    CITEXT,
  "title"   TEXT,
  "votes"   INT4 DEFAULT '0'
);

DROP INDEX IF EXISTS index_on_threads_slug;

CREATE UNIQUE INDEX index_on_threads_slug ON threads (slug) WHERE slug != '';

DROP INDEX IF EXISTS index_on_threads_forum;

CREATE INDEX index_on_threads_forum ON threads (forum);

DROP INDEX IF EXISTS index_on_threads_slug_and_created;

CREATE INDEX index_on_threads_slug_and_created ON threads (slug, created);


CREATE TABLE "votes" (
  "voice"  INT2,
  "user"   CITEXT,
  "thread" BIGINT
);

DROP INDEX IF EXISTS index_on_votes_user_and_thread;

CREATE UNIQUE INDEX index_on_votes_user_and_thread ON votes (thread, "user");


CREATE TABLE IF NOT EXISTS "members" (
  forum  CITEXT,
  author CITEXT
);

DROP INDEX IF EXISTS index_on_members_forum_and_author;

CREATE UNIQUE INDEX index_on_members_forum_and_author ON members (forum, author);