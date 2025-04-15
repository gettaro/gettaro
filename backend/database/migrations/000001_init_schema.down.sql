-- Drop tables in reverse order of creation to handle foreign key constraints
DROP TABLE IF EXISTS pm_tickets;
DROP TABLE IF EXISTS project_management_accounts;
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS pr_comments;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS source_control_accounts;
DROP TABLE IF EXISTS integration_configs;
DROP TABLE IF EXISTS direct_reports;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS user_titles;
DROP TABLE IF EXISTS titles;
DROP TABLE IF EXISTS auth_providers;
DROP TABLE IF EXISTS user_organizations;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
