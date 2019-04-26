-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS roles (
  id SERIAL NOT NULL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,  
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL
);  

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS roles;
