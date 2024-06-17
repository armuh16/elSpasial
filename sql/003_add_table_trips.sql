-- +goose Up
create table trips (
    id              bigserial primary key,
    origin          varchar(255) not null,
    destination     varchar(255) not null,
    user_id         int not null,
    price           float not null,
    updated_at      timestamptz default now(),
    created_at      timestamptz default now(),
    deleted_at      timestamptz default null,
    foreign key (user_id) references users (id)
);

-- +goose Down
drop table trips;