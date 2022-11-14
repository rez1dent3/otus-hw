alter table events
    add is_sent bool not null;

create index is_sent_ind
    on events (is_sent);

