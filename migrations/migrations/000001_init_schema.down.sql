DROP TRIGGER IF EXISTS update_secrets_updated_at ON secrets;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS secrets;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";