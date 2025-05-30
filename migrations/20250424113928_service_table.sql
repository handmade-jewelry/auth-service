-- +goose Up
create table if not exists service (
    id serial primary key,
    name varchar(255) not null unique,
    is_active boolean default false,
    host varchar(255) not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

-- +goose Down
drop table if exists service;
