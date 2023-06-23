-- +migrate Up
DROP INDEX companies_id_idx;
CREATE INDEX CONCURRENTLY companies_id_idx ON companies (id);

DROP INDEX companies_title_idx;
CREATE INDEX CONCURRENTLY companies_title_idx ON companies (title);

DROP INDEX clients_id_idx;
CREATE INDEX CONCURRENTLY clients_id_idx ON clients (id);

-- +migrate Down
DROP INDEX companies_id_idx;
DROP INDEX companies_title_idx;