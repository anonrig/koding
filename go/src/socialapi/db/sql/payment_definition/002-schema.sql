-- SET ROLE social;

CREATE SCHEMA payment;

GRANT usage ON SCHEMA payment to social;

-- append sitemap schema
SELECT set_config('search_path', current_setting('search_path') || ',payment', false);
