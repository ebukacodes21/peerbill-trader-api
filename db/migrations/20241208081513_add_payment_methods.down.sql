ALTER TABLE "payment_methods" 
  DROP CONSTRAINT IF EXISTS "crypto_fiat_payment_method_key";

ALTER TABLE "payment_methods" 
  DROP CONSTRAINT IF EXISTS "fk_username";

DROP TABLE IF EXISTS "payment_methods";