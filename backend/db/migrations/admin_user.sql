BEGIN;
INSERT INTO users
VALUES (DEFAULT, 'admin@dev.com', 'Admin', 'User', '+33600000000', '1900-01-01');
INSERT INTO accounts
VALUES (1, '__ADMIN_PASSWORD_HASH__', DEFAULT, 'admin');
COMMIT;
