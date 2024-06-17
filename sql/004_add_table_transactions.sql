-- +goose Up
create table transactions (
    id              bigserial primary key,
    driver_id       int not null,
    user_id         int not null,
    trip           json not null,
    grand_total     float not null,
    status          int not null default 1,
    updated_at      timestamptz default now(),
    created_at      timestamptz default now(),
    deleted_at      timestamptz default null,
--     foreign key  (driver_id) references users (id),
    foreign key  (user_id) references users (id)
);

-- +goose Down
drop table transactions;