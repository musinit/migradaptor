-- +migrate Up
DROP INDEX companies_id_idx;
CREATE INDEX CONCURRENTLY companies_id_idx ON companies (id);