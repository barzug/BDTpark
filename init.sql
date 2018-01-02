DROP TABLE IF EXISTS "forums";

DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "posts";

DROP TABLE IF EXISTS "threads";

DROP TABLE IF EXISTS "votes";

DROP TABLE IF EXISTS "members";

CREATE TABLE "forums" (
  "fID"     SERIAL NOT NULL,
  "posts"   INT8   NOT NULL DEFAULT '0',
  "slug"    CITEXT NOT NULL,
  "threads" INT4   NOT NULL DEFAULT '0',
  "title"   TEXT   NOT NULL,
  "author"  TEXT   NOT NULL,
  CONSTRAINT forums_pk PRIMARY KEY ("fID")
) WITH (
OIDS = FALSE
);

DROP INDEX IF EXISTS index_on_forums_slug;

CREATE UNIQUE INDEX  index_on_forums_slug
  ON forums (slug);


CREATE TABLE "users" (
  "uID"      SERIAL NOT NULL,
  "email"    CITEXT NOT NULL,
  "nickname" CITEXT NOT NULL,
  "fullname" TEXT,
  "about"    TEXT,
  CONSTRAINT users_pk PRIMARY KEY ("uID")
) WITH (
OIDS = FALSE
);

DROP INDEX IF EXISTS index_on_users_email;

CREATE UNIQUE INDEX index_on_users_email
  ON users (email);

DROP INDEX IF EXISTS index_on_users_nickname;

CREATE UNIQUE INDEX index_on_users_nickname
  ON users (nickname);


CREATE TABLE "posts" (
  "pID"      SERIAL  NOT NULL,
  "author"   TEXT    NOT NULL,
  "created"  TIMESTAMP WITH TIME ZONE DEFAULT now(),
  "forum"    TEXT    NOT NULL,
  "isEdited" BOOLEAN NOT NULL         DEFAULT 'false',
  "message"  TEXT,
  "parent"   INT8                     DEFAULT '0',
  "thread"   BIGINT  NOT NULL,
  "path"     INT8 [],
  CONSTRAINT posts_pk PRIMARY KEY ("pID")
) WITH (
OIDS = FALSE
);

DROP INDEX IF EXISTS index_on_posts_thread;

CREATE INDEX index_on_posts_thread
  ON posts (thread);

DROP INDEX IF EXISTS index_on_posts_parent;

CREATE INDEX index_on_posts_parent
  ON posts (parent);

DROP INDEX IF EXISTS index_on_posts_path;

CREATE INDEX index_on_posts_path
  ON posts (path);


CREATE TABLE "threads" (
  "tID"     SERIAL NOT NULL,
  "author"  TEXT   NOT NULL,
  "created" TIMESTAMP WITH TIME ZONE DEFAULT now(),
  "forum"   CITEXT NOT NULL,
  "message" TEXT   NOT NULL,
  "slug"    CITEXT,
  "title"   TEXT   NOT NULL,
  "votes"   INT4   NOT NULL          DEFAULT '0',
  CONSTRAINT threads_pk PRIMARY KEY ("tID")
) WITH (
OIDS = FALSE
);

DROP INDEX IF EXISTS threads_slug_key;

CREATE UNIQUE INDEX threads_slug_key ON threads (slug) WHERE slug != '';

DROP INDEX IF EXISTS threads_forum_key;

CREATE INDEX threads_forum_key ON threads (forum);


CREATE TABLE "votes" (
  "voice"  INT2,
  "user"   CITEXT,
  "thread" BIGINT
) WITH (
OIDS = FALSE
);

DROP INDEX IF EXISTS index_on_votes_user_and_thread;

CREATE UNIQUE INDEX index_on_votes_user_and_thread ON votes (thread, "user");


CREATE TABLE IF NOT EXISTS "members" (
  forum  CITEXT,
  author CITEXT
) WITH (
  OIDS = FALSE
);

DROP INDEX IF EXISTS index_on_members_forum_and_author;

CREATE UNIQUE INDEX index_on_members_forum_and_author ON members (forum, author);
