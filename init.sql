ALTER TABLE "posts" DROP CONSTRAINT IF EXISTS "posts_fk0";

ALTER TABLE "posts" DROP CONSTRAINT IF EXISTS "posts_fk1";

ALTER TABLE "threads" DROP CONSTRAINT IF EXISTS "threads_fk0";

ALTER TABLE "threads" DROP CONSTRAINT IF EXISTS "threads_fk1";

ALTER TABLE "votes" DROP CONSTRAINT IF EXISTS "votes_fk0";

ALTER TABLE "votes" DROP CONSTRAINT IF EXISTS "votes_fk1";

DROP TABLE IF EXISTS "forums";

DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "posts";

DROP TABLE IF EXISTS "threads";

DROP TABLE IF EXISTS "votes";

CREATE TABLE "forums" (
  "fID" serial NOT NULL,
  "posts" int8 NOT NULL DEFAULT '0',
  "slug" CITEXT NOT NULL UNIQUE,
  "threads" int4 NOT NULL DEFAULT '0',
  "title" TEXT NOT NULL,
  "author" TEXT NOT NULL,
  CONSTRAINT forums_pk PRIMARY KEY ("fID")
) WITH (
OIDS=FALSE
);


CREATE TABLE "users" (
  "uID" serial NOT NULL,
  "email" CITEXT NOT NULL UNIQUE,
  "nickname" CITEXT NOT NULL UNIQUE,
  "fullname" TEXT,
  "about" TEXT,
  CONSTRAINT users_pk PRIMARY KEY ("uID")
) WITH (
OIDS=FALSE
);


CREATE TABLE "posts" (
  "pID" serial NOT NULL,
  "author" bigint NOT NULL,
  "created" DATE NOT NULL,
  "forum" TEXT NOT NULL,
  "isEdited" BOOLEAN NOT NULL DEFAULT 'false',
  "message" TEXT,
  "parent" int8 DEFAULT '0',
  "thread" bigint NOT NULL,
  CONSTRAINT posts_pk PRIMARY KEY ("pID")
) WITH (
OIDS=FALSE
);



CREATE TABLE "threads" (
  "tID" bigint NOT NULL,
  "author" bigint NOT NULL,
  "created" DATE NOT NULL,
  "forum" bigint NOT NULL,
  "description" TEXT NOT NULL,
  "slug" TEXT NOT NULL UNIQUE,
  "title" TEXT NOT NULL,
  "votes" int4 NOT NULL DEFAULT '0',
  CONSTRAINT threads_pk PRIMARY KEY ("tID")
) WITH (
OIDS=FALSE
);



CREATE TABLE "votes" (
  "voice" int2,
  "user" bigint NOT NULL,
  "thread" bigint NOT NULL
) WITH (
OIDS=FALSE
);



ALTER TABLE "posts" ADD CONSTRAINT "posts_fk0" FOREIGN KEY ("author") REFERENCES "users"("uID");
ALTER TABLE "posts" ADD CONSTRAINT "posts_fk1" FOREIGN KEY ("thread") REFERENCES "threads"("tID");

ALTER TABLE "threads" ADD CONSTRAINT "threads_fk0" FOREIGN KEY ("author") REFERENCES "users"("uID");
ALTER TABLE "threads" ADD CONSTRAINT "threads_fk1" FOREIGN KEY ("forum") REFERENCES "forums"("fID");

ALTER TABLE "votes" ADD CONSTRAINT "votes_fk0" FOREIGN KEY ("user") REFERENCES "users"("uID");
ALTER TABLE "votes" ADD CONSTRAINT "votes_fk1" FOREIGN KEY ("thread") REFERENCES "threads"("tID");