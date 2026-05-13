INSERT INTO tenants (name, slug, status)
VALUES ('Moodora Local', 'moodora-local', 'active')
ON CONFLICT (lower(slug)) DO UPDATE
SET name = EXCLUDED.name,
    status = EXCLUDED.status,
    updated_at = now();

INSERT INTO users (id, tenant_id, email, display_name, status)
SELECT :'dev_user_id'::uuid, tenants.id, 'dev@moodora.local', 'Moodora Dev User', 'active'
FROM tenants
WHERE tenants.slug = 'moodora-local'
ON CONFLICT (tenant_id, lower(email)) DO UPDATE
SET display_name = EXCLUDED.display_name,
    status = EXCLUDED.status,
    updated_at = now();

INSERT INTO user_profiles (user_id, timezone, locale)
VALUES (:'dev_user_id'::uuid, 'Asia/Bangkok', 'th')
ON CONFLICT (user_id) DO UPDATE
SET timezone = EXCLUDED.timezone,
    locale = EXCLUDED.locale,
    updated_at = now();
