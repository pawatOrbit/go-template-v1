-- Drop trigger first
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop the users table
DROP TABLE IF EXISTS users;

-- Note: We don't drop the trigger function as it might be used by other tables
-- If you're sure no other tables use it, uncomment the line below:
-- DROP FUNCTION IF EXISTS update_updated_at_column();