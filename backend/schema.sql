CREATE TABLE "memos" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "user_id" uuid NOT NULL, "title" text NOT NULL, "text" text NOT NULL, PRIMARY KEY ("id"));
CREATE TABLE "rooms" ("id" uuid NOT NULL, "name" text NOT NULL, "created_at" timestamp NOT NULL DEFAULT now(), PRIMARY KEY ("id"), UNIQUE ("name"));
CREATE TABLE "messages" ("id" uuid NOT NULL, "user_id" uuid NOT NULL, "room_id" uuid NOT NULL, "text" text NOT NULL, "created_at" timestamp NOT NULL DEFAULT now(), PRIMARY KEY ("id"), FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE);
