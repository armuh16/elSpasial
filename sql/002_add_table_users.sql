-- +goose Up
create table users (
    id        bigserial primary key,
    name      varchar(255) not null,
    password  varchar(255) not null,
    email     varchar(255) unique not null,
    address   varchar(255) not null,
    role_id      int not null,
    updated_at timestamptz default now(),
    created_at timestamptz default now(),
    deleted_at timestamptz default null,
    foreign key (role_id) references roles (id)
);

-- +goose Down
drop table users;