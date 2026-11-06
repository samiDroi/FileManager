-- +goose Up

-- Registration creates a user before any other user exists, so the author of
-- the row cannot be known. These audit columns are not used to authorize
-- access and must accept NULL during bootstrap.
ALTER TABLE public."user" ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE public."user" ALTER COLUMN updated_by DROP NOT NULL;

-- +goose Down

UPDATE public."user" SET created_by = 0 WHERE created_by IS NULL;
UPDATE public."user" SET updated_by = 0 WHERE updated_by IS NULL;

ALTER TABLE public."user" ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE public."user" ALTER COLUMN updated_by SET NOT NULL;
