ALTER TABLE users  add IF NOT EXISTS created_at TIMESTAMP;
ALTER TABLE users  add IF NOT EXISTS updated_at TIMESTAMP;

ALTER TABLE roles  add IF NOT EXISTS status BOOLEAN DEFAULT true;
ALTER TABLE roles  add IF NOT EXISTS created_at TIMESTAMP;
ALTER TABLE roles  add IF NOT EXISTS updated_at TIMESTAMP;

ALTER TABLE permision_role  add IF NOT EXISTS status BOOLEAN DEFAULT true;
ALTER TABLE user_role  add IF NOT EXISTS status BOOLEAN DEFAULT true;

DELETE FROM permissions;

INSERT INTO permissions (id, name) VALUES
                                       (1, 'Users'),
                                       (2, 'Roles'),
                                       (3, 'Services'),
                                       (4, 'Service fees'),
                                       (5, 'Pricing fees'),
                                       (6, 'Subscriptions') ON CONFLICT DO NOTHING;
