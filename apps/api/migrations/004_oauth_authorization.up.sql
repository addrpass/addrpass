-- OAuth2 authorization codes (short-lived, single-use)
CREATE TABLE authorization_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT UNIQUE NOT NULL,
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address_id UUID NOT NULL REFERENCES addresses(id) ON DELETE CASCADE,
    scope TEXT NOT NULL DEFAULT 'full' CHECK (scope IN ('full', 'delivery', 'zone', 'verify')),
    redirect_uri TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_authorization_codes_code ON authorization_codes(code);

-- OAuth apps (registered by businesses for the consent flow)
CREATE TABLE oauth_apps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    logo_url TEXT NOT NULL DEFAULT '',
    redirect_uris TEXT[] NOT NULL DEFAULT '{}',
    client_id TEXT UNIQUE NOT NULL,
    client_secret_hash TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_oauth_apps_client_id ON oauth_apps(client_id);
CREATE INDEX idx_oauth_apps_business_id ON oauth_apps(business_id);
