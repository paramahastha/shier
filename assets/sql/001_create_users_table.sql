-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS users (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT NOT NULL,
  "email" TEXT NOT NULL,
  "password" TEXT NOT NULL,
  "confirm" TEXT NOT NULL,
  "role" TEXT NOT NULL,
  "created_at" TIMESTAMP NULL,
  "updated_at" TIMESTAMP NULL
);  

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS users;
