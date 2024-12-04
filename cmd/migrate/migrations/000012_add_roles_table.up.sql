CREATE TABLE IF NOT EXISTS roles(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
    roles(name, description, level)
    VALUES (
        'user', 
        'a user can create posts and comments',
        1
);

INSERT INTO
    roles(name, description, level)
    VALUES (
        'moderator', 
        'a moderator can update other users posts',
        2
);

INSERT INTO
    roles(name, description, level)
    VALUES (
        'admin', 
        'an admin can update or delete user posts',
        3
);