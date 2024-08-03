CREATE TABLE "memos" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "title" text NOT NULL,
    "text" text NOT NULL,
    PRIMARY KEY ("id")
);