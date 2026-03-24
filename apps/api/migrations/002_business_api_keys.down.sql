ALTER TABLE access_logs DROP COLUMN IF EXISTS business_name;
ALTER TABLE access_logs DROP COLUMN IF EXISTS business_id;
ALTER TABLE shares DROP COLUMN IF EXISTS scope;
DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF EXISTS webhooks;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS businesses;
