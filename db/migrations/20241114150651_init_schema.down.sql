ALTER TABLE "trade_pairs" 
  DROP CONSTRAINT IF EXISTS "crypto_fiat_trader_key";

ALTER TABLE "trade_pairs" 
  DROP CONSTRAINT IF EXISTS "fk_username";

DROP TABLE IF EXISTS "trade_pairs";

DROP TABLE IF EXISTS "traders";
