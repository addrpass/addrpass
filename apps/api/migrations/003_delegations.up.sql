-- Delegation chains: allow share owners to delegate access to businesses
-- Customer → E-commerce → Delivery Company → Driver
CREATE TABLE delegations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    share_id UUID NOT NULL REFERENCES shares(id) ON DELETE CASCADE,
    from_business_id UUID REFERENCES businesses(id) ON DELETE CASCADE,
    to_business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    scope TEXT NOT NULL DEFAULT 'full' CHECK (scope IN ('full', 'delivery', 'zone', 'verify')),
    expires_at TIMESTAMPTZ,
    max_accesses INT,
    access_count INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_delegations_share_id ON delegations(share_id);
CREATE INDEX idx_delegations_to_business_id ON delegations(to_business_id);

-- User-to-business delegation (first link in chain: user shares token with a business)
ALTER TABLE shares ADD COLUMN delegated_to_business_id UUID REFERENCES businesses(id);

-- Labels generated for shares
CREATE TABLE labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    share_id UUID NOT NULL REFERENCES shares(id) ON DELETE CASCADE,
    business_id UUID REFERENCES businesses(id),
    reference_code TEXT UNIQUE NOT NULL,
    zone_code TEXT NOT NULL DEFAULT '',
    format TEXT NOT NULL DEFAULT 'pdf' CHECK (format IN ('pdf', 'png', 'zpl')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_labels_reference_code ON labels(reference_code);
CREATE INDEX idx_labels_share_id ON labels(share_id);
