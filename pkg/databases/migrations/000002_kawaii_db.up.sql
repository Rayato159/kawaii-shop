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
    "role_id",
    "created_at"
)
VALUES
(),
();

COMMIT;