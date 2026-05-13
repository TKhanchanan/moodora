DROP TRIGGER IF EXISTS coin_transactions_prevent_delete ON coin_transactions;
DROP TRIGGER IF EXISTS coin_transactions_prevent_update ON coin_transactions;
DROP FUNCTION IF EXISTS prevent_coin_transaction_mutation();

DROP TRIGGER IF EXISTS api_keys_set_updated_at ON api_keys;
DROP TRIGGER IF EXISTS check_ins_set_updated_at ON check_ins;
DROP TRIGGER IF EXISTS wallets_set_updated_at ON wallets;
DROP TRIGGER IF EXISTS user_profiles_set_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS users_set_updated_at ON users;
DROP TRIGGER IF EXISTS tenants_set_updated_at ON tenants;
DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS api_usage_logs;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS check_ins;
DROP TABLE IF EXISTS coin_transactions;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;
