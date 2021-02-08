SELECT 'CREATE DATABASE libreria_test'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'libreria_test')\gexec
