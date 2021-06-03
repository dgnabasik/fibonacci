CREATE TABLE IF NOT EXISTS public.fibonacci (
    id integer NOT NULL DEFAULT 0,
    fibvalue numeric(308,0),
   	CONSTRAINT fibonacci_pkey PRIMARY KEY (id)
)
WITH (OIDS=FALSE) TABLESPACE pg_default;
