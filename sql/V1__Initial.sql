CREATE TABLE IF NOT EXISTS roles (
                                     id              SERIAL PRIMARY KEY,
                                     name           VARCHAR(256) NOT NULL
    );

CREATE TABLE IF NOT EXISTS orgs (
                                    id              SERIAL PRIMARY KEY,
                                    name           VARCHAR(256) NOT NULL,
    type           INTEGER NOT NULL
    );

CREATE TABLE IF NOT EXISTS contracts (
                                         id              SERIAL PRIMARY KEY,
                                         supply_vendor_id           INTEGER NOT NULL,
                                         start_date           TIMESTAMP WITHOUT TIME ZONE,
                                         end_date           TIMESTAMP WITHOUT TIME ZONE,
                                         base_amount	DECIMAL(12, 2),
    actual_amount	DECIMAL(12, 2),
    code           VARCHAR(256) NOT NULL
    );

CREATE TABLE IF NOT EXISTS permissions (
                                           id              SERIAL PRIMARY KEY,
                                           name           VARCHAR(256) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS permision_role (
                                              role_id              SERIAL NOT NULL,
                                              permission_id        SERIAL NOT NULL
);

INSERT INTO roles (id, name) VALUES
                                 (1, 'Planner Team'),
                                 (2, 'Project Contractor'),
                                 (3, 'Supply Vendor') ON CONFLICT DO NOTHING;

INSERT INTO orgs (id, name, type) VALUES
                                      (3, '3M', 3),
                                      (4, 'SanMar', 3),
                                      (1, 'SPG Company', 1),
                                      (2, 'TPP Technologies', 2) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS users (
                                     id              SERIAL PRIMARY KEY,
                                     email           VARCHAR(1024) UNIQUE NOT NULL,
    full_name       VARCHAR(1024) NOT NULL,
    is_admin        BOOLEAN DEFAULT FALSE,
    org_id	INTEGER NOT NULL,
    password	VARCHAR(1024) NOT NULL
    );

CREATE TABLE IF NOT EXISTS user_role (
                                         user_id              SERIAL NOT NULL,
                                         role_id        SERIAL NOT NULL
);

INSERT INTO contracts (id, supply_vendor_id, start_date, end_date, base_amount, actual_amount, code) VALUES
                                                                                                         (12, 4, '2022-07-01 00:00:00', '2022-10-01 00:00:00', 250, 250, 'SANMAR006'),
                                                                                                         (9, 4, '2022-07-25 00:00:00', '2022-08-25 00:00:00', 90, 90, 'SANMAR003'),
                                                                                                         (10, 4, '2022-07-08 00:00:00', '2022-07-14 00:00:00', 125, 125, 'SANMAR004'),
                                                                                                         (11, 4, '2022-07-12 00:00:00', '2022-08-01 00:00:00', 175, 175, 'SANMAR005'),
                                                                                                         (7, 4, '2022-07-06 00:00:00', '2022-09-25 00:00:00', 50, 50, 'SANMAR001'),
                                                                                                         (8, 4, '2022-07-15 00:00:00', '2022-10-01 00:00:00', 80, 80, 'SANMAR002'),
                                                                                                         (3, 3, '2022-07-11 00:00:00', '2022-11-11 00:00:00', 80, 80, '3M003'),
                                                                                                         (4, 3, '2022-07-13 00:00:00', '2022-12-01 00:00:00', 150, 150, '3M004'),
                                                                                                         (5, 3, '2022-07-05 00:00:00', '2022-08-12 00:00:00', 200, 200, '3M005'),
                                                                                                         (6, 3, '2022-07-22 00:00:00', '2022-09-01 00:00:00', 170, 170, '3M006'),
                                                                                                         (1, 3, '2022-07-04 00:00:00', '2022-08-04 00:00:00', 100, 100, '3M001'),
                                                                                                         (2, 3, '2022-07-01 00:00:00', '2022-08-01 00:00:00', 50, 50, '3M002') ON CONFLICT DO NOTHING;



INSERT INTO permissions (id, name) VALUES
                                       (1, 'VIEW_ALL_CONTRACT_LIST'),
                                       (2, 'VIEW_CONTRACT_LIST'),
                                       (3, 'VIEW_ALL_WITHDRAW_REQUEST'),
                                       (4, 'VIEW_WITHDRAW_REQUEST'),
                                       (5, 'CREATE_WITHDRAW_REQUEST'),
                                       (6, 'UPDATE_WITHDRAW_REQUEST'),
                                       (7, 'CANCEL_WITHDRAW_REQUEST'),
                                       (8, 'UPDATE_WITHDRAW_REQUEST_STATUS_TO_READY_TO_COLLECT'),
                                       (9, 'UPDATE_WITHDRAW_REQUEST_STATUS_TO_COLLECTED') ON CONFLICT DO NOTHING;


INSERT INTO permision_role (role_id, permission_id) VALUES
                                                        (1, 1),
                                                        (1, 3),
                                                        (1, 5),
                                                        (1, 6),
                                                        (1, 7),
                                                        (2, 1),
                                                        (2, 4),
                                                        (2, 9),
                                                        (3, 2),
                                                        (3, 4),
                                                        (3, 8) ON CONFLICT DO NOTHING;


INSERT INTO users (id, email, full_name, is_admin, org_id, password) VALUES
                                                                         (2, 'projectcontractor.spgtpp@gmail.com', 'Tom Hank', false, 2, '$2a$10$xnZmDRV/IlylHykjD4dxUeBPqX3OrorIw6oIprfW7XaZkjcJLJ.hq'),
                                                                         (3, 'supplyvendor.spgtpp@gmail.com', 'Tom Hiddleston', false, 3, '$2a$10$xnZmDRV/IlylHykjD4dxUeBPqX3OrorIw6oIprfW7XaZkjcJLJ.hq'),
                                                                         (4, 'admin.spgtpp@gmail.com', 'Tom Holland', true, 1, '$2a$12$9vtukhuSn.CxkX77qLmd0uM91arwaznZHbq4dxJ72paOj0Pb6n5b2'),
                                                                         (5, 'supplyvendor2.spgtpp@gmail.com', 'Tom Hardy', false, 4, '$2a$10$xnZmDRV/IlylHykjD4dxUeBPqX3OrorIw6oIprfW7XaZkjcJLJ.hq'),
                                                                         (1, 'planner.spgtpp@gmail.com', 'Tom Cruise', false, 1, '$2a$10$xnZmDRV/IlylHykjD4dxUeBPqX3OrorIw6oIprfW7XaZkjcJLJ.hq') ON CONFLICT DO NOTHING;


INSERT INTO user_role (user_id, role_id) VALUES
                                             (1, 1),
                                             (2, 2),
                                             (3, 3),
                                             (5, 3) ON CONFLICT DO NOTHING;

SELECT pg_catalog.setval('contracts_id_seq', 1, false);


SELECT pg_catalog.setval('orgs_id_seq', 1, false);


SELECT pg_catalog.setval('permissions_id_seq', 1, false);


SELECT pg_catalog.setval('roles_id_seq', 1, false);


SELECT pg_catalog.setval('users_id_seq', 1, false);