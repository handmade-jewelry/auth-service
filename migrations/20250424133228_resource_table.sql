-- +goose Up
create table if not exists resource (
    id serial primary key,
    service_id int not null,
    public_path text not null,
    service_path text not null,
    is_active bool default false,
    roles jsonb,
    method varchar(10),
    scheme varchar(10) not null,
    check_access_token bool default true,
    check_roles bool default true,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp
);

create unique index unique_active_resource
    on resource(service_id, public_path, method)
    where deleted_at is null;

-- +goose Down
drop table if exists resource;

