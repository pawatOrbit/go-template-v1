# Database Migrations

This directory contains database migration files for the Go API Template project. Migrations are managed using [golang-migrate](https://github.com/golang-migrate/migrate).

## 📁 Migration File Structure

Migration files follow this naming convention:
```
{version}_{description}.{direction}.sql

Examples:
001_create_users_table.up.sql    # Migration up
001_create_users_table.down.sql  # Migration rollback
002_create_api_keys_table.up.sql
002_create_api_keys_table.down.sql
```

## 🚀 Quick Commands

### Running Migrations

```bash
# Run all pending migrations
make migrate
# or
make migrate-up

# Check migration status
make migrate-status

# Show current version
make migrate-version
```

### Rollback Migrations

```bash
# Rollback last migration
make migrate-down

# Rollback all migrations (DANGEROUS!)
make migrate-down-all

# Rollback specific number of steps
go run main.go migrate down --steps=3
```

### Creating New Migrations

```bash
# Create a new migration
make migrate-create name=add_user_preferences_table

# This creates:
# migrations/20240101120000_add_user_preferences_table.up.sql
# migrations/20240101120000_add_user_preferences_table.down.sql
```

### Advanced Commands

```bash
# Force database to specific version (use with caution)
make migrate-force version=2

# Check detailed migration status
go run main.go migrate status
```

## 📝 Migration Best Practices

### 1. **Always Create Both UP and DOWN**
- Every `.up.sql` file should have a corresponding `.down.sql` file
- Down migrations should safely reverse the changes made by up migrations

### 2. **Use Transactions When Possible**
```sql
BEGIN;

-- Your migration statements here
CREATE TABLE example (...);
CREATE INDEX idx_example ON example (...);

COMMIT;
```

### 3. **Handle Existing Data Carefully**
```sql
-- Good: Check if column exists before adding
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

-- Good: Use default values for NOT NULL columns
ALTER TABLE users 
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE 
DEFAULT CURRENT_TIMESTAMP NOT NULL;
```

### 4. **Create Indexes Concurrently in Production**
```sql
-- For large tables in production
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email 
ON users(email);
```

### 5. **Use Proper Data Types**
```sql
-- Good: Use appropriate data types
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Example Migrations

### Creating a Table
**Up Migration** (`001_create_users_table.up.sql`):
```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

**Down Migration** (`001_create_users_table.down.sql`):
```sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS users;
```

### Adding a Column
**Up Migration** (`004_add_user_phone.up.sql`):
```sql
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone) 
WHERE phone IS NOT NULL;
```

**Down Migration** (`004_add_user_phone.down.sql`):
```sql
DROP INDEX IF EXISTS idx_users_phone;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### Data Migration
**Up Migration** (`005_migrate_user_roles.up.sql`):
```sql
-- Add new roles column
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS roles TEXT[] DEFAULT ARRAY['user'];

-- Migrate existing data
UPDATE users 
SET roles = ARRAY['admin'] 
WHERE email LIKE '%@company.com';

-- Remove old role column if exists
ALTER TABLE users DROP COLUMN IF EXISTS role;
```

**Down Migration** (`005_migrate_user_roles.down.sql`):
```sql
-- Add back old role column
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'user';

-- Migrate data back
UPDATE users 
SET role = CASE 
    WHEN 'admin' = ANY(roles) THEN 'admin'
    ELSE 'user'
END;

-- Remove new roles column
ALTER TABLE users DROP COLUMN IF EXISTS roles;
```

## ⚠️ Safety Guidelines

### Development Environment
- Always test migrations on a copy of production data
- Run migrations in a transaction when possible
- Keep backups before running destructive migrations

### Production Environment
- **Always backup the database before running migrations**
- Test migrations on staging environment first
- Consider maintenance windows for large migrations
- Use `CREATE INDEX CONCURRENTLY` for large tables
- Monitor migration progress for long-running operations

### Emergency Rollback
If a migration fails:

1. **Check the error message**:
   ```bash
   make migrate-status
   ```

2. **If database is in dirty state**, fix manually or force to last known good version:
   ```bash
   make migrate-force version=N
   ```

3. **Rollback if needed**:
   ```bash
   make migrate-down
   ```

## 🔍 Troubleshooting

### Common Issues

**"database is locked"**
- Another migration process is running
- Stop other processes and retry

**"dirty database version"**
- A migration failed partway through
- Check database state and use `migrate force` if necessary

**"no change"**
- All migrations are already applied
- Check `make migrate-status` for current state

**Connection errors**
- Verify database is running: `make db-up`
- Check configuration in `config/config.local.yaml`

### Getting Help

```bash
# Show help for migration commands
go run main.go migrate --help

# Show help for specific subcommand
go run main.go migrate up --help
```

## 📚 Additional Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Migration Best Practices](https://www.prisma.io/blog/database-migrations-best-practices)