CREATE EXTENSION pgcrypto;

CREATE TABLE greetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    greeting TEXT NOT NULL DEFAULT '',
    created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

create table tidedata (
    id        serial not null,
    lat       double precision,
    long      double precision,
    extremes  jsonb,
    heights   jsonb,
    timestamp bigint
);