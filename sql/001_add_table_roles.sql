-- +goose Up
create table roles (
    id            bigserial primary key,
    name          varchar(50),
    updated_at    timestamptz default now(),
    created_at    timestamptz default now()
);

insert into roles(id, name) values (1, 'Driver');
insert into roles(id, name) values (2, 'User');

-- +goose Down
drop table roles;