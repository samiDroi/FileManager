-- +goose Up

ALTER TABLE public."sessions" ADD COLUMN IF NOT EXISTS family_id VARCHAR(64);
ALTER TABLE public."sessions" ADD COLUMN IF NOT EXISTS replaced_by_id BIGINT REFERENCES public."sessions" (id) ON DELETE SET NULL;

UPDATE public."sessions" SET family_id = id::text WHERE family_id IS NULL OR family_id = '';
ALTER TABLE public."sessions" ALTER COLUMN family_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_sessions_family_id ON public."sessions" (family_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON public."sessions" (expires_at);

-- +goose Down

DROP INDEX IF EXISTS public.idx_sessions_expires_at;
DROP INDEX IF EXISTS public.idx_sessions_family_id;

ALTER TABLE public."sessions" DROP COLUMN IF EXISTS replaced_by_id;
ALTER TABLE public."sessions" DROP COLUMN IF EXISTS family_id;
