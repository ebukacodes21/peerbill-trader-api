CREATE TABLE "traders" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "country" varchar NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "profile_pic" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trade_pairs" (
  "id" bigserial PRIMARY KEY,
  "username" VARCHAR NOT NULL,
  "crypto" varchar NOT NULL,
  "fiat" varchar NOT NULL,
  "buy_rate" double precision NOT NULL,
  "sell_rate" double precision NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "trade_pairs" ADD CONSTRAINT "fk_username" FOREIGN KEY ("username") REFERENCES "traders" ("username");

-- CREATE UNIQUE INDEX ON "trade_pairs" ("crypto", "fiat", "username")
ALTER TABLE "trade_pairs" ADD CONSTRAINT "base_quote_trader_key" UNIQUE ("crypto", "fiat", "username");

