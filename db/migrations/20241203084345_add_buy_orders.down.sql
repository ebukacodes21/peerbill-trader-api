ALTER TABLE "buy_orders" 
  DROP CONSTRAINT IF EXISTS "fk_username";

DROP TABLE IF EXISTS "buy_orders";