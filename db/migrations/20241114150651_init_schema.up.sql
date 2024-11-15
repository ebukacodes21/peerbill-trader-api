CREATE TABLE "traders" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "country" varchar NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "status" varchar NOT NULL DEFAULT 'disabled',
  "role" varchar NOT NULL DEFAULT 'user',
  "profile_pic" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trade_pairs" (
  "id" bigserial PRIMARY KEY,
  "trader_id" bigint NOT NULL,
  "base_asset" varchar NOT NULL,
  "quote_asset" varchar NOT NULL,
  "buy_rate" double precision NOT NULL,
  "sell_rate" double precision NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "trade_pairs" ADD CONSTRAINT "fk_trader_id" FOREIGN KEY ("trader_id") REFERENCES "traders" ("id");

-- CREATE UNIQUE INDEX ON "trade_pairs" ("base_asset", "quote_asset", "trader_id")
ALTER TABLE "trade_pairs" ADD CONSTRAINT "base_quote_trader_key" UNIQUE ("base_asset", "quote_asset", "trader_id");

