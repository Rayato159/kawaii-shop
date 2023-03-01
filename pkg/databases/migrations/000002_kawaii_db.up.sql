BEGIN;

INSERT INTO "roles" (
    "title"
)
VALUES
('user'),
('admin');

INSERT INTO "users" (
    "username",
    "email",
    "token",
    "role_id"
)
VALUES
('test_user', 'test@kawaii.com', ''),
();

COMMIT;