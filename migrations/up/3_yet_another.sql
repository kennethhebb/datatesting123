create table yet_another (
    id        serial not null,
    extremes  jsonb,
    heights   jsonb,
    timestamp bigint
);