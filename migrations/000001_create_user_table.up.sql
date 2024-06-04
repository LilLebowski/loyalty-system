CREATE TABLE IF NOT EXISTS "user"
(
    id            uuid      default gen_random_uuid() PRIMARY KEY,
    login         varchar(100) not null,
    password_hash varchar(100) not null,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_login UNIQUE (login)
);
