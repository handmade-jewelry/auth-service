-- +goose Up
create table if not exists resource (
    id serial primary key,
    service_id int not null,
    path text not null,
    is_active bool default false,
    roles text,
    method varchar(10),
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

-- +goose Down
drop table if exists resource;

