CREATE TABLE IF NOT EXISTS "balance"
(
    id        uuid default gen_random_uuid() PRIMARY KEY,
    user_id   uuid NOT NULL,
    current   DECIMAL(10, 2),
    withdrawn DECIMAL(10, 2),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);
