BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Domains configured for sending
CREATE TABLE IF NOT EXISTS domains (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  domain VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  verified_at TIMESTAMPTZ,
  is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_domains_domain ON domains(domain) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_domains_verified ON domains(verified_at) WHERE verified_at IS NOT NULL;

-- API keys to use domains
CREATE TABLE IF NOT EXISTS api_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  domain_id UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
  key_hash VARCHAR(64) NOT NULL UNIQUE,
  key_prefix VARCHAR(16) NOT NULL,
  name VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_used_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ,
  is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_api_keys_domain ON api_keys(domain_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash) WHERE is_active = true;

-- Email send history
CREATE TABLE IF NOT EXISTS emails (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  domain_id UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
  api_key_id UUID NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
  message_id TEXT NOT NULL,
  from_address VARCHAR(255) NOT NULL,
  to_address VARCHAR(255) NOT NULL,
  subject TEXT,
  smtp_response_code SMALLINT,
  smtp_response_text TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_emails_domain_crated_id ON emails(domain_id, created_at DESC, id DESC);

COMMIT;
