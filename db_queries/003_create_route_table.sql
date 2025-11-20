CREATE TABLE stops (
    id          SERIAL PRIMARY KEY,
    route_id    INT REFERENCES routes(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    stop_order  INT NOT NULL,  -- 1, 2, 3,... in sequence along the route
    geom        geometry(Point, 4326)
);
