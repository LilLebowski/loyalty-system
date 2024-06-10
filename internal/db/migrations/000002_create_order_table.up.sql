CREATE TYPE "order_status" AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
    );

CREATE TABLE IF NOT EXISTS "order"
(
    id          uuid      default gen_random_uuid() PRIMARY KEY,
    user_id     uuid         NOT NULL,
    number      varchar(100) NOT NULL,
    status      order_status,
    accrual     DECIMAL(10, 2),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);
