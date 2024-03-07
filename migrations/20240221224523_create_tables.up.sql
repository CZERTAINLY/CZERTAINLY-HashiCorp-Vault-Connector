create schema hvault


create table hvault.authority_instances
(
    id                 serial,
    uuid               varchar(255) not null,
    name               varchar(255),
    url                varchar(255),
    credential_type    varchar(255),
    role_id    varchar(255),
    role_secret    varchar(255),
    jwt    varchar(255),
    attributes         text,
    primary key (id)
);

create table hvault.certificates
(
    id                  serial,
    serial_number       varchar not null,
    uuid                varchar not null,
    base64_content       varchar null default null,
    meta                text null default null,
    primary key (id)
);

create table hvault.discoveries
(
    id      serial,
    uuid    varchar not null,
    name    varchar not null,
    status  varchar not null,
    meta    text null default null,
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
