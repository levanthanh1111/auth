ALTER TABLE users  add IF NOT EXISTS status BOOLEAN DEFAULT true;
ALTER TABLE users  add IF NOT EXISTS last_login TIMESTAMP WITHOUT TIME ZONE;
