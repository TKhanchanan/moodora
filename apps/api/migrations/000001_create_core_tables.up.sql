CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tenants_slug_not_empty CHECK (length(trim(slug)) > 0),
    CONSTRAINT tenants_status_valid CHECK (status IN ('active', 'suspended', 'disabled'))
);

CREATE UNIQUE INDEX tenants_slug_unique_idx ON tenants (lower(slug));
CREATE INDEX tenants_status_idx ON tenants (status);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE RESTRICT,
    email TEXT NOT NULL,
    display_name TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_email_not_empty CHECK (length(trim(email)) > 0),
    CONSTRAINT users_status_valid CHECK (status IN ('active', 'suspended', 'deleted'))
);

CREATE UNIQUE INDEX users_tenant_email_unique_idx ON users (tenant_id, lower(email));
CREATE INDEX users_tenant_id_idx ON users (tenant_id);
CREATE INDEX users_status_idx ON users (status);
CREATE INDEX users_created_at_idx ON users (created_at);

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    birth_date DATE,
    timezone TEXT NOT NULL DEFAULT 'Asia/Bangkok',
    locale TEXT NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT user_profiles_timezone_not_empty CHECK (length(trim(timezone)) > 0),
    CONSTRAINT user_profiles_locale_not_empty CHECK (length(trim(locale)) > 0)
);

CREATE INDEX user_profiles_user_id_idx ON user_profiles (user_id);

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users (id) ON DELETE RESTRICT,
    coin_balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT wallets_coin_balance_non_negative CHECK (coin_balance >= 0)
);

CREATE INDEX wallets_user_id_idx ON wallets (user_id);

CREATE TABLE coin_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets (id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    transaction_type TEXT NOT NULL,
    amount BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    reason TEXT NOT NULL,
    idempotency_key TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT coin_transactions_amount_not_zero CHECK (amount <> 0),
    CONSTRAINT coin_transactions_balance_after_non_negative CHECK (balance_after >= 0),
    CONSTRAINT coin_transactions_type_valid CHECK (transaction_type IN ('grant', 'spend', 'adjustment', 'refund')),
    CONSTRAINT coin_transactions_reason_not_empty CHECK (length(trim(reason)) > 0)
);

CREATE INDEX coin_transactions_wallet_created_at_idx ON coin_transactions (wallet_id, created_at DESC);
CREATE INDEX coin_transactions_user_created_at_idx ON coin_transactions (user_id, created_at DESC);
CREATE UNIQUE INDEX coin_transactions_wallet_idempotency_unique_idx
    ON coin_transactions (wallet_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

CREATE TABLE check_ins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    local_date DATE NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'Asia/Bangkok',
    reward_coins BIGINT NOT NULL DEFAULT 0,
    coin_transaction_id UUID REFERENCES coin_transactions (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT check_ins_reward_coins_non_negative CHECK (reward_coins >= 0),
    CONSTRAINT check_ins_timezone_not_empty CHECK (length(trim(timezone)) > 0)
);

CREATE UNIQUE INDEX check_ins_user_local_date_unique_idx ON check_ins (user_id, local_date);
CREATE INDEX check_ins_user_created_at_idx ON check_ins (user_id, created_at DESC);
CREATE INDEX check_ins_local_date_idx ON check_ins (local_date);

CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    key_prefix TEXT NOT NULL,
    key_hash TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT api_keys_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT api_keys_prefix_not_empty CHECK (length(trim(key_prefix)) > 0),
    CONSTRAINT api_keys_hash_not_empty CHECK (length(trim(key_hash)) > 0),
    CONSTRAINT api_keys_status_valid CHECK (status IN ('active', 'revoked', 'expired'))
);

CREATE UNIQUE INDEX api_keys_key_hash_unique_idx ON api_keys (key_hash);
CREATE INDEX api_keys_tenant_id_idx ON api_keys (tenant_id);
CREATE INDEX api_keys_key_prefix_idx ON api_keys (key_prefix);
CREATE INDEX api_keys_status_idx ON api_keys (status);

COMMENT ON COLUMN api_keys.key_hash IS 'Hash of the API key. Raw API keys must never be stored.';

CREATE TABLE api_usage_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    api_key_id UUID REFERENCES api_keys (id) ON DELETE SET NULL,
    route TEXT NOT NULL,
    method TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    request_count INTEGER NOT NULL DEFAULT 1,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT api_usage_logs_route_not_empty CHECK (length(trim(route)) > 0),
    CONSTRAINT api_usage_logs_method_not_empty CHECK (length(trim(method)) > 0),
    CONSTRAINT api_usage_logs_status_code_valid CHECK (status_code BETWEEN 100 AND 599),
    CONSTRAINT api_usage_logs_request_count_positive CHECK (request_count > 0)
);

CREATE INDEX api_usage_logs_tenant_requested_at_idx ON api_usage_logs (tenant_id, requested_at DESC);
CREATE INDEX api_usage_logs_api_key_requested_at_idx ON api_usage_logs (api_key_id, requested_at DESC);
CREATE INDEX api_usage_logs_route_requested_at_idx ON api_usage_logs (route, requested_at DESC);

CREATE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tenants_set_updated_at
BEFORE UPDATE ON tenants
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER user_profiles_set_updated_at
BEFORE UPDATE ON user_profiles
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER wallets_set_updated_at
BEFORE UPDATE ON wallets
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER check_ins_set_updated_at
BEFORE UPDATE ON check_ins
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER api_keys_set_updated_at
BEFORE UPDATE ON api_keys
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE FUNCTION prevent_coin_transaction_mutation()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'coin_transactions are immutable';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER coin_transactions_prevent_update
BEFORE UPDATE ON coin_transactions
FOR EACH ROW
EXECUTE FUNCTION prevent_coin_transaction_mutation();

CREATE TRIGGER coin_transactions_prevent_delete
BEFORE DELETE ON coin_transactions
FOR EACH ROW
EXECUTE FUNCTION prevent_coin_transaction_mutation();
