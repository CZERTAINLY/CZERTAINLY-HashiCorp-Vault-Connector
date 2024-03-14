create schema hvault


create table hvault.authority_instances
(
    id                 serial,
    uuid               varchar(255) not null unique ,
    name               varchar(255) not null unique ,
    url                varchar(255) not null ,
    credential_type    varchar(255) not null,
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
    uuid                varchar not null unique ,
    base64_content       varchar null default null,
    meta                text null default null,
    primary key (id)
);

CREATE INDEX index_certificates_serial_number ON hvault.certificates (serial_number);
CREATE INDEX index_certificates_uuid ON hvault.certificates (uuid);

create table hvault.discoveries
(
    id      serial,
    uuid    varchar not null unique ,
    name    varchar not null unique,
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
