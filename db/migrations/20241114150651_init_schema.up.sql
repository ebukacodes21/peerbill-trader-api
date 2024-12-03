-- Create the "traders" table first, as other tables depend on it
CREATE TABLE "traders" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_verified" bool NOT NULL DEFAULT false,
  "verification_code" varchar UNIQUE NOT NULL,
  "country" varchar NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "profile_pic" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create the "trade_pairs" table
CREATE TABLE "trade_pairs" (
  "id" bigserial PRIMARY KEY,
  "username" VARCHAR NOT NULL,
  "crypto" varchar NOT NULL,
  "fiat" varchar NOT NULL,
  "buy_rate" double precision NOT NULL,
  "sell_rate" double precision NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Add the foreign key constraint for the "trade_pairs" table to reference the "traders" table
ALTER TABLE "trade_pairs" 
  ADD CONSTRAINT "fk_username" FOREIGN KEY ("username") REFERENCES "traders" ("username");

-- Create a unique constraint to ensure no duplicate trade pairs for the same trader
ALTER TABLE "trade_pairs" 
  ADD CONSTRAINT "crypto_fiat_trader_key" UNIQUE ("crypto", "fiat", "username");