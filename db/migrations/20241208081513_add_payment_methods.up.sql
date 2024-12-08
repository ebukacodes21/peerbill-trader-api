CREATE TABLE "payment_methods" (
    "id" bigserial PRIMARY KEY,
    "trade_pair_id" bigint NOT NULL,
    "username" varchar NOT NULL,
    "wallet_address" varchar NOT NULL,
    "crypto" varchar NOT NULL,
    "fiat" varchar NOT NULL,
    "bank_name" varchar NOT NULL,
    "account_number" varchar NOT NULL,
    "account_holder" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT(now())
);

-- Add foreign key constraint for the "payment_methods" table referencing the "trade_pairs" table
ALTER TABLE "payment_methods" 
  ADD CONSTRAINT "fk_username" FOREIGN KEY ("trade_pair_id") REFERENCES "trade_pairs" ("id");

  -- Create a unique constraint to ensure no duplicate payment method for the same trade pair
ALTER TABLE "payment_methods" 
  ADD CONSTRAINT "crypto_fiat_payment_method_key" UNIQUE ("crypto", "fiat", "trade_pair_id");