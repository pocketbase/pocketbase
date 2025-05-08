--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

-- Enable built-in pgcrypto extension to use gen_random_bytes function
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Adding "nocase" collation to be compatible with SQLite's default "nocase" collation
CREATE COLLATION IF NOT EXISTS "nocase" (
  provider = icu,          -- Specify ICU as the provider
  locale = 'und-u-ks-level2', -- Undetermined locale, Unicode extension (-u-), collation strength (ks) level 2 (level2)
  deterministic = false    -- Case-insensitive collations are typically non-deterministic
);

-- Alias [hex] to encode(..., 'hex')
CREATE OR REPLACE FUNCTION hex(data bytea)
RETURNS text
LANGUAGE SQL
IMMUTABLE
AS $$
SELECT encode(data, 'hex')
$$;

-- Alias [randomblob] to gen_random_bytes(...)
CREATE OR REPLACE FUNCTION randomblob(length integer)
RETURNS bytea
LANGUAGE SQL
IMMUTABLE
AS $$
SELECT gen_random_bytes(length)
$$;

-- Create the uuid_generate_v7 function
create or replace function uuid_generate_v7()
    returns uuid
    as $$
    begin
    -- use random v4 uuid as starting point (which has the same variant we need)
    -- then overlay timestamp
    -- then set version 7 by flipping the 2 and 1 bit in the version 4 string
    return encode(
        set_bit(
        set_bit(
            overlay(uuid_send(gen_random_uuid())
                    placing substring(int8send(floor(extract(epoch from clock_timestamp()) * 1000)::bigint) from 3)
                    from 1 for 6
            ),
            52, 1
        ),
        53, 1
        ),
        'hex')::uuid;
    end
    $$
    language plpgsql
    volatile;

-- Create json_valid function
CREATE OR REPLACE FUNCTION json_valid(text) RETURNS boolean AS $$
	BEGIN
		PERFORM $1::jsonb;
		RETURN TRUE;
	EXCEPTION WHEN others THEN
		RETURN FALSE;
	END;
	$$ LANGUAGE plpgsql IMMUTABLE;

--
-- Name: _logs; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7() NOT NULL,
    level bigint DEFAULT '0'::bigint,
    message text DEFAULT '""'::text,
    data json DEFAULT '"{}"'::json,
    created TIMESTAMP DEFAULT NOW() NOT NULL,
    updated TIMESTAMP DEFAULT NOW() NOT NULL
);


ALTER TABLE public._logs OWNER TO "user";

--
-- Name: _migrations; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._migrations (
    file text,
    applied bigint
);


ALTER TABLE public._migrations OWNER TO "user";

--
-- Data for Name: _logs; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._logs (id, level, message, data, created, updated) FROM stdin;
\.


--
-- Data for Name: _migrations; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._migrations (file, applied) FROM stdin;
1640988000_init.go	1679427415788852
1660821103_add_user_ip_column.go	1679427415792425
1677760279_uppsercase_method.go	1679427415792706
1699187560_logs_generalization.go	1700504854831333
\.


--
-- PostgreSQL database dump complete
--

