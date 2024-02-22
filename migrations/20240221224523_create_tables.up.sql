create schema hvault

create sequence hvault.authority_instances_id_seq start 1 increment 1;

create table hvault.authority_instances
(
    id                 int8         not null,
    uuid               varchar(255) not null,
    name               varchar(255),
    url                varchar(255),
    credential_uuid    varchar(255),
    credential_data    text,
    attributes         json,
    primary key (id)
);

create sequence hvault.certificates_id_seq start 1 increment 1;
create sequence hvault.discoveries_id_seq start 1 increment 1;

create table hvault.certificates
(
    id              bigint not null,
    serial_id       varchar not null,
    uuid            varchar not null,
    base64content   varchar null default null,
    meta            json null default null,
    primary key (id)
);

create table hvault.discoveries
(
    id      bigint not null,
    uuid    varchar not null,
    name    varchar not null,
    status  varchar not null,
    meta    json null default null,
    primary key (id)
);

create table hvault.discovery_certificates
(
    certificate_id bigint not null,
    discovery_id   bigint not null,
    primary key (certificate_id, discovery_id),
    foreign key (certificate_id) references hvault.certificates(id),
    foreign key (discovery_id) references hvault.discoveries(id)
);
