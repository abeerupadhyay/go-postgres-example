begin;

-- function to set `created_at` field as current time
create or replace function set_created_at_column()
returns trigger as $$
begin
    NEW.created_at = now();
    return NEW;
end;
$$ language plpgsql;

-- function to set `updated_at` field as current time
create or replace function set_updated_at_column()
returns trigger as $$
begin
    NEW.updated_at = now();
    return NEW;
end;
$$ language plpgsql;

-- new table `book`
create table book (

    -- auto generated fields
    id           bigserial primary key,
    created_at   timestamp with time zone not null,
    updated_at   timestamp with time zone null,

    -- user defined fields
    isbn          varchar (17) not null,
    title         varchar (256) not null,
    author        varchar (256) not null,
    publish_year  char (4) not null,
    rating        decimal
);

-- create indexes on book table
create unique index unique_book_isbn on book (isbn);

-- trigger to set `created_at` field on every insert in book table
create trigger book_set_created_at
    before insert on book
    for each row execute procedure set_created_at_column();

-- trigger to set `updated_at` field on every update in book table
create trigger book_set_updated_at
    before update on book
    for each row execute procedure set_updated_at_column();

commit;
