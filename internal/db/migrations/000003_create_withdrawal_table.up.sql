CREATE TABLE IF NOT EXISTS "withdrawal"
(
    id                uuid      default gen_random_uuid() PRIMARY KEY,
    user_id           uuid         NOT NULL,
    external_order_id varchar(100) NOT NULL,
    sum               DECIMAL(10, 2),
    processed_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);