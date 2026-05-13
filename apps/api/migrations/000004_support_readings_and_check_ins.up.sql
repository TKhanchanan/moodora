ALTER TABLE tarot_readings
ADD COLUMN question TEXT;

ALTER TABLE coin_transactions
DROP CONSTRAINT coin_transactions_type_valid;

ALTER TABLE coin_transactions
ADD CONSTRAINT coin_transactions_type_valid
CHECK (transaction_type IN ('grant', 'spend', 'adjustment', 'refund', 'check_in'));
