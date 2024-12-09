CREATE TABLE "orders" (
    "id" bigserial PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "escrow_address" VARCHAR NOT NULL,
    "user_address" VARCHAR NOT NULL,
    "order_type" VARCHAR NOT NULL,
    "crypto" VARCHAR NOT NULL,
    "fiat" VARCHAR NOT NULL,
    "crypto_amount" double precision NOT NULL,
    "fiat_amount" double precision NOT NULL,
    "rate" double precision NOT NULL,
    "is_accepted" BOOL NOT NULL DEFAULT false,
    "is_completed" BOOL NOT NULL DEFAULT false,
    "is_rejected" BOOL NOT NULL DEFAULT false,
    "is_expired" BOOL NOT NULL DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "duration" TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30 minutes'
);

-- Add foreign key constraint for the "orders" table referencing the "traders" table
ALTER TABLE "orders" 
  ADD CONSTRAINT "fk_username" FOREIGN KEY ("username") REFERENCES "traders" ("username");