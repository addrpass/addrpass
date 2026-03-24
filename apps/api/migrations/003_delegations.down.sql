DROP TABLE IF EXISTS labels;
ALTER TABLE shares DROP COLUMN IF EXISTS delegated_to_business_id;
DROP TABLE IF EXISTS delegations;
