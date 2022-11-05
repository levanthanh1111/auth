DELETE FROM permision_role;

INSERT INTO permision_role (role_id, permission_id) VALUES
                                                        (1, 1),
                                                        (1, 3),
                                                        (1, 5),
                                                        (1, 6),
                                                        (2, 1),
                                                        (2, 4),
                                                        (3, 2),
                                                        (3, 4) ON CONFLICT DO NOTHING;