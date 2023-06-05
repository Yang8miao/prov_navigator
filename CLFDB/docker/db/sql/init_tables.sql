drop table if exists `r_log_tag`;
drop table if exists `log`;
drop table if exists `tag`;
drop table if exists `dns`;

create table if not exists log
(
    _id      int unsigned auto_increment
        primary key,
    _time    timestamp(3)           not null,
    log_raw  text                   null,
    log_type varchar(64) default '' not null
);

create table if not exists clf_db.tag
(
    _id    int unsigned auto_increment
        primary key,
    _key   varchar(64)   not null,
    _value varchar(512)  not null,
    _type  int default 0 not null comment '0 - normal, 1 - group',
    constraint tag_kvt_index
        unique (_key, _value, _type)
);

create table clf_db.r_log_tag
(
    _id    int unsigned auto_increment
        primary key,
    log_id int unsigned not null,
    tag_id int unsigned not null,
    constraint r_log_tag_id_index
        unique (log_id, tag_id),
    constraint r_log_tag_log_id_fk
        foreign key (log_id) references clf_db.log (_id)
            on delete cascade,
    constraint r_log_tag_tag_id_fk
        foreign key (tag_id) references clf_db.tag (_id)
            on delete cascade
);

create table if not exists clf_db.dns
(
    _id      int unsigned auto_increment
    primary key,
    _domain  varchar(64)            not null,
    _ip      varchar(32)            not null,
    constraint r_dns
    unique (_domain, _ip)
    );
