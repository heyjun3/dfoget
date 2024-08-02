-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";

-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';

CREATE TABLE IF NOT EXISTS "memos" (
    id UUID NOT NULL PRIMARY KEY,
    title VARCHAR NOT NULL,
    text TEXT NOT NULL
);