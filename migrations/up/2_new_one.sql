create table new_one (
    id        serial not null,
    extremes  jsonb,
    heights   jsonb,
    timestamp bigint
);