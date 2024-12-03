CREATE TABLE "buy_orders" (
    "id" bigserial PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "wallet_address" varchar NOT NULL,
    "crypto" varchar NOT NULL,
    "fiat" varchar NOT NULL,
    "crypto_amount" double precision NOT NULL,
    "fiat_amount" double precision NOT NULL,
    "rate" double precision NOT NULL,
    "is_accepted" bool NOT NULL DEFAULT false,
    "is_completed" bool NOT NULL DEFAULT false,
    "is_rejected" bool NOT NULL DEFAULT false,
    "is_expired" bool NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT(now()),
    "duration" timestamptz NOT NULL
);

-- Add foreign key constraint for the "buy_orders" table referencing the "traders" table
ALTER TABLE "buy_orders" 
  ADD CONSTRAINT "fk_username" FOREIGN KEY ("username") REFERENCES "traders" ("username");