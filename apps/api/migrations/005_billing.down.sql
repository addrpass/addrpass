DROP TABLE IF EXISTS billing_events;
DROP TABLE IF EXISTS usage_records;
ALTER TABLE users DROP COLUMN IF EXISTS plan_resolution_limit;
ALTER TABLE users DROP COLUMN IF EXISTS stripe_subscription_id;
ALTER TABLE users DROP COLUMN IF EXISTS stripe_customer_id;
ALTER TABLE users DROP COLUMN IF EXISTS plan;
