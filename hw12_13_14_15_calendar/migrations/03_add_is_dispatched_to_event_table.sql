alter table events
    rename column is_sent to in_queue;

alter table events
    add is_dispatched bool not null;

create index is_dispatched_ind
    on events (is_dispatched);

