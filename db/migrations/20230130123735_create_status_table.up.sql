CREATE TABLE IF NOT EXISTS "links"(
    "id" BIGINT SERIAL PRIMARY KEY,
    "link" TEXT NOT NULL,
    "status" VARCHAR(255) CHECK
        ("status" IN('UP','DOWN')) NOT NULL
);
