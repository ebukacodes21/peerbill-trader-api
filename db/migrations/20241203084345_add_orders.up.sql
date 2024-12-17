CREATE TABLE "orders" (
    "id" bigserial PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "escrow_address" VARCHAR,
    "user_address" VARCHAR NOT NULL,
    "order_type" VARCHAR NOT NULL,
    "crypto" VARCHAR NOT NULL,
    "fiat" VARCHAR NOT NULL,
    "crypto_amount" double precision NOT NULL,
    "fiat_amount" double precision NOT NULL,
    "bank_name" varchar,
    "account_number" varchar,
    "account_holder" varchar,
    "rate" double precision NOT NULL,
    "is_accepted" BOOL NOT NULL DEFAULT false,
    "is_completed" BOOL NOT NULL DEFAULT false,
    "is_rejected" BOOL NOT NULL DEFAULT false,
    "is_received" BOOL NOT NULL DEFAULT false,
    "is_expired" BOOL NOT NULL DEFAULT false,
    "duration" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add foreign key constraint for the "orders" table referencing the "traders" table
ALTER TABLE "orders" 
  ADD CONSTRAINT "fk_username" FOREIGN KEY ("username") REFERENCES "traders" ("username");