
DROP TABLE IF EXISTS "forums";

DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "posts";

DROP TABLE IF EXISTS "threads";

DROP TABLE IF EXISTS "votes";

DROP TABLE IF EXISTS "members";

CREATE TABLE "forums" (
  "fID"     SERIAL NOT NULL,
  "posts"   INT8   NOT NULL DEFAULT '0',
  "slug"    CITEXT NOT NULL UNIQUE,
  "threads" INT4   NOT NULL DEFAULT '0',
  "title"   TEXT   NOT NULL,
  "author"  TEXT   NOT NULL,
  CONSTRAINT forums_pk PRIMARY KEY ("fID")
) WITH (
OIDS = FALSE
);


CREATE TABLE "users" (
  "uID"      SERIAL NOT NULL,
  "email"    CITEXT NOT NULL UNIQUE,
  "nickname" CITEXT NOT NULL UNIQUE,
  "fullname" TEXT,
  "about"    TEXT,
  CONSTRAINT users_pk PRIMARY KEY ("uID")
) WITH (
OIDS = FALSE
);


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

CREATE UNIQUE INDEX IF NOT EXISTS threads_slug_key ON threads (slug) WHERE slug != '';



CREATE TABLE "votes" (
  "voice"  INT2,
  "user"   CITEXT NOT NULL,
  "thread" BIGINT NOT NULL,
  UNIQUE ("user", "thread")
) WITH (
OIDS = FALSE
);


CREATE TABLE IF NOT EXISTS "members" (
  forum  CITEXT,
author CITEXT,
UNIQUE ("forum", "author")
) WITH (
  OIDS = FALSE
);
