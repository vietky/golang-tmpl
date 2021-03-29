CREATE DATABASE cfs;

use cfs;

CREATE TABLE IF NOT EXISTS call_for_services (
    agency_id VARCHAR(100),
    event_id VARCHAR(100),
    event_number VARCHAR(100),
    event_time TIMESTAMP,
    dispatch_time TIMESTAMP,
    responder VARCHAR(100),
    PRIMARY KEY (event_id)
);


CREATE TABLE IF NOT EXISTS agencies (
    id VARCHAR(100) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(100) PRIMARY KEY,
    agency_id VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS responders (
    id VARCHAR(100) PRIMARY KEY,
    agency_id VARCHAR(100)
);

ALTER TABLE call_for_services ADD FOREIGN KEY (agency_id) REFERENCES agencies(id);
ALTER TABLE users ADD FOREIGN KEY (agency_id) REFERENCES agencies(id);
ALTER TABLE responders ADD FOREIGN KEY (agency_id) REFERENCES agencies(id);
