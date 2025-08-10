-- Drop trigger first
DROP TRIGGER IF EXISTS update_api_keys_updated_at ON api_keys;

-- Drop the api_keys table
DROP TABLE IF EXISTS api_keys;