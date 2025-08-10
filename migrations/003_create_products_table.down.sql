-- Drop trigger first
DROP TRIGGER IF EXISTS update_products_updated_at ON products;

-- Drop the products table
DROP TABLE IF EXISTS products;