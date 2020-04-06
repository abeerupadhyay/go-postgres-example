begin;

-- drop table
drop table if exists book;

-- drop functions
drop function if exists public.set_created_at_column;
drop function if exists public.set_updated_at_column;

commit;
