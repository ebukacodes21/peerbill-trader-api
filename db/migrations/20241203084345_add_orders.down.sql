ALTER TABLE "orders" 
  DROP CONSTRAINT IF EXISTS "fk_username";

DROP TABLE IF EXISTS "orders";