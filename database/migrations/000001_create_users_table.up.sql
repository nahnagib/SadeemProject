CREATE TABLE users (
    id            uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    img           VARCHAR(255),
    name          VARCHAR(255) NOT NULL,
    phone         VARCHAR(255) NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password      VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
