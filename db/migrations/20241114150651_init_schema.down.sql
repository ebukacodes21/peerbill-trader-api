ALTER TABLE "trade_pairs" 
  DROP CONSTRAINT IF EXISTS "base_quote_trader_key";

ALTER TABLE "trade_pairs" 
  DROP CONSTRAINT IF EXISTS "fk_trader_id";

DROP TABLE IF EXISTS "trade_pairs";

DROP TABLE IF EXISTS "traders";
