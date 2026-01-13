-- verify_env.sql
-- Verifies:
-- 1) current database
-- 2) current user
-- 3) role 'docker' exists & can login
-- 4) database 'order' exists & is owned by 'docker'
-- 5) TCP host/port of this session

\echo
\echo ==== CHECK: current database & user ====
SELECT current_database()  AS current_database,
       current_user        AS current_user;

\echo
\echo ==== CHECK: role 'docker' exists & can login ====
SELECT rolname, rolcanlogin, rolsuper
FROM pg_roles
WHERE rolname = 'docker';

\echo
\echo ==== CHECK: database 'order' exists & owner ====
SELECT d.datname AS db_name,
       r.rolname AS owner
FROM pg_database d
JOIN pg_roles r ON r.oid = d.datdba
WHERE d.datname = 'order';
    
\echo
\echo ==== CHECK: server TCP address & port for this session ====
SELECT inet_server_addr() AS server_host,
       inet_server_port() AS server_port;

\echo
\echo ==== OPTIONAL: show search_path & server_version ====
SHOW search_path;
SHOW server_version;
