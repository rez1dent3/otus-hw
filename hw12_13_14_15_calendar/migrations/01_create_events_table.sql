create table events
(
    id          uuid        not null
        constraint id
            primary key,
    title       varchar(30) not null,
    description text,
    start_at    timestamptz not null,
    end_at      timestamptz not null,
    user_id     uuid        not null,
    remind_for  time,
    created_at  timestamptz not null,
    updated_at  timestamptz not null
);

create index user_ind
    on events (user_id);

