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
-- Name: "_authOrigins"; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public."_authOrigins" (
    "collectionRef" text DEFAULT ''::text,
    created text DEFAULT ''::text,
    fingerprint text DEFAULT ''::text,
    id text DEFAULT '''r''||lower(hex(randomblob(7)))'::text NOT NULL,
    "recordRef" text DEFAULT ''::text,
    updated text DEFAULT ''::text
);


ALTER TABLE public."_authOrigins" OWNER TO "user";

--
-- Name: _collections; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._collections (
    id text NOT NULL,
    system boolean DEFAULT false,
    type text DEFAULT '"base"'::text,
    name text,
    fields jsonb DEFAULT '"[]"'::jsonb,
    "listRule" text,
    "viewRule" text,
    "createRule" text,
    "updateRule" text,
    "deleteRule" text,
    options jsonb DEFAULT '"{}"'::jsonb,
    created text DEFAULT '""'::text,
    updated text DEFAULT '""'::text,
    indexes jsonb DEFAULT '"[]"'::jsonb
);


ALTER TABLE public._collections OWNER TO "user";

--
-- Name: "_externalAuths"; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public."_externalAuths" (
    "collectionRef" text DEFAULT ''::text,
    created text DEFAULT ''::text,
    id text DEFAULT '''r''||lower(hex(randomblob(7)))'::text NOT NULL,
    provider text DEFAULT ''::text,
    "providerId" text DEFAULT ''::text,
    "recordRef" text DEFAULT ''::text,
    updated text DEFAULT ''::text
);


ALTER TABLE public."_externalAuths" OWNER TO "user";

--
-- Name: _mfas; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._mfas (
    "collectionRef" text DEFAULT ''::text,
    created text DEFAULT ''::text,
    id text DEFAULT '''r''||lower(hex(randomblob(7)))'::text NOT NULL,
    method text DEFAULT ''::text,
    "recordRef" text DEFAULT ''::text,
    updated text DEFAULT ''::text
);


ALTER TABLE public._mfas OWNER TO "user";

--
-- Name: _migrations; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._migrations (
    file text NOT NULL,
    applied bigint
);


ALTER TABLE public._migrations OWNER TO "user";

--
-- Name: _otps; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._otps (
    "collectionRef" text DEFAULT ''::text,
    created text DEFAULT ''::text,
    id text DEFAULT '''r''||lower(hex(randomblob(7)))'::text NOT NULL,
    password text DEFAULT ''::text,
    "recordRef" text DEFAULT ''::text,
    updated text DEFAULT ''::text,
    "sentTo" text DEFAULT ''::text
);


ALTER TABLE public._otps OWNER TO "user";

--
-- Name: _params; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._params (
    created TEXT DEFAULT ''::text NOT NULL,
    id TEXT DEFAULT '''r''||lower(hex(randomblob(7)))'::text NOT NULL,
    updated TEXT DEFAULT '' NOT NULL,
    value TEXT DEFAULT NULL -- Use TEXT because encrypted values are not valid JSON.
);

ALTER TABLE public._params OWNER TO "user";

--
-- Name: _superusers; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public._superusers (
    created text DEFAULT ''::text,
    email text DEFAULT ''::text,
    "emailVisibility" boolean DEFAULT false,
    id text DEFAULT '''r''||lower(hex(randomblob(7)))'::text NOT NULL,
    password text DEFAULT ''::text,
    "tokenKey" text DEFAULT ''::text,
    updated text DEFAULT ''::text,
    verified boolean DEFAULT false
);


ALTER TABLE public._superusers OWNER TO "user";

--
-- Name: clients; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.clients (
    created text DEFAULT ''::text,
    email text DEFAULT ''::text,
    "emailVisibility" boolean DEFAULT false,
    id text NOT NULL,
    name text DEFAULT ''::text,
    password text,
    "tokenKey" text,
    updated text DEFAULT ''::text,
    username text,
    verified boolean DEFAULT false
);


ALTER TABLE public.clients OWNER TO "user";

--
-- Name: demo1; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.demo1 (
    created text DEFAULT ''::text,
    id text NOT NULL,
    updated text DEFAULT ''::text,
    text text DEFAULT ''::text,
    bool boolean DEFAULT false,
    url text DEFAULT ''::text,
    select_one text DEFAULT ''::text,
    file_one text DEFAULT ''::text,
    file_many text DEFAULT ''::text,
    number real DEFAULT '0'::real,
    email text DEFAULT ''::text,
    datetime text DEFAULT ''::text,
    "json" jsonb,
    rel_one text DEFAULT ''::text,
    select_many text DEFAULT ''::text,
    rel_many jsonb DEFAULT '[]'::jsonb,
    point jsonb DEFAULT '{"lat":0,"lon":0}'::jsonb
);


ALTER TABLE public.demo1 OWNER TO "user";

--
-- Name: demo2; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.demo2 (
    created text DEFAULT ''::text,
    id text NOT NULL,
    title text DEFAULT ''::text,
    updated text DEFAULT ''::text,
    active boolean DEFAULT false
);


ALTER TABLE public.demo2 OWNER TO "user";

--
-- Name: demo3; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.demo3 (
    created text DEFAULT ''::text,
    id text NOT NULL,
    title text DEFAULT ''::text,
    updated text DEFAULT ''::text,
    files text DEFAULT ''::text
);


ALTER TABLE public.demo3 OWNER TO "user";

--
-- Name: demo4; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.demo4 (
    created text DEFAULT ''::text,
    id text NOT NULL,
    title text DEFAULT ''::text,
    updated text DEFAULT ''::text,
    rel_one_no_cascade text DEFAULT ''::text,
    rel_one_no_cascade_required text DEFAULT ''::text,
    rel_one_cascade text DEFAULT ''::text,
    self_rel_one text DEFAULT ''::text,
    "json_array" jsonb,
    "json_object" jsonb,
    rel_one_unique text DEFAULT ''::text,
    rel_many_no_cascade jsonb DEFAULT '[]'::jsonb,
    rel_many_no_cascade_required jsonb DEFAULT '[]'::jsonb,
    rel_many_cascade jsonb DEFAULT '[]'::jsonb,
    rel_many_unique jsonb DEFAULT '[]'::jsonb,
    self_rel_many jsonb DEFAULT '[]'::jsonb
);


ALTER TABLE public.demo4 OWNER TO "user";

--
-- Name: demo5; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.demo5 (
    created text DEFAULT ''::text,
    id text NOT NULL,
    rel_many jsonb DEFAULT '[]'::jsonb,
    rel_one text DEFAULT ''::text,
    select_many text DEFAULT ''::text,
    select_one text DEFAULT ''::text,
    updated text DEFAULT ''::text,
    total real DEFAULT '0'::real,
    file text DEFAULT ''::text
);


ALTER TABLE public.demo5 OWNER TO "user";

--
-- Name: nologin; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.nologin (
    created text DEFAULT ''::text,
    email text DEFAULT ''::text,
    "emailVisibility" boolean DEFAULT false,
    id text NOT NULL,
    name text DEFAULT ''::text,
    password text,
    "tokenKey" text,
    updated text DEFAULT ''::text,
    username text,
    verified boolean DEFAULT false
);


ALTER TABLE public.nologin OWNER TO "user";

--
-- Name: sqlite_stat1; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.sqlite_stat1 (
    tbl text,
    idx text,
    stat text
);


ALTER TABLE public.sqlite_stat1 OWNER TO "user";

--
-- Name: sqlite_stat4; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.sqlite_stat4 (
    tbl text,
    idx text,
    neq text,
    nlt text,
    ndlt text,
    sample text
);


ALTER TABLE public.sqlite_stat4 OWNER TO "user";

--
-- Name: users; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.users (
    avatar text DEFAULT ''::text,
    created text DEFAULT ''::text,
    email text DEFAULT ''::text,
    "emailVisibility" boolean DEFAULT false,
    id text NOT NULL,
    name text DEFAULT ''::text,
    password text,
    "tokenKey" text,
    updated text DEFAULT ''::text,
    username text,
    verified boolean DEFAULT false,
    file text DEFAULT ''::text,
    rel text DEFAULT ''::text
);


ALTER TABLE public.users OWNER TO "user";

--
-- Data for Name: "_authOrigins"; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public."_authOrigins" ("collectionRef", created, fingerprint, id, "recordRef", updated) FROM stdin;
pbc_3142635823	2024-07-13 09:24:15.442Z	6afbfe481c31c08c55a746cccb88ece0	5f29jy38bf5zm3f	sywbhecnh46rhm0	2024-11-04 16:41:03.658Z
pbc_3142635823	2024-07-26 12:10:47.102Z	22bbbcbed36e25321f384ccf99f60057	dmy260k6ksjr4ib	sbmbsdb40jyxf7h	2024-11-04 16:41:03.656Z
pbc_3142635823	2024-07-26 12:11:38.697Z	6afbfe481c31c08c55a746cccb88ece0	ic55o70g4f8pcl4	sbmbsdb40jyxf7h	2024-11-04 16:41:03.654Z
pbc_3142635823	2024-07-26 12:12:17.972Z	dc879cfc889d0f1c1f3258d6f3a828fe	5798yh833k6w6w0	sbmbsdb40jyxf7h	2024-11-04 16:41:03.650Z
v851q4r790rhknl	2024-07-26 12:22:37.681Z	22bbbcbed36e25321f384ccf99f60057	9r2j0m74260ur8i	gk390qegs4y47wn	2024-07-26 12:22:37.681Z
\.


--
-- Data for Name: _collections; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._collections (id, system, type, name, fields, "listRule", "viewRule", "createRule", "updateRule", "deleteRule", options, created, updated, indexes) FROM stdin;
_pb_users_auth_	f	auth	users	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"cost":10,"hidden":true,"id":"_pbf_auth_password_","max":0,"min":8,"name":"password","pattern":"","presentable":false,"required":true,"system":true,"type":"password"},{"autogeneratePattern":"[a-zA-Z0-9_]{50}","hidden":true,"id":"_pbf_auth_tokenKey_","max":60,"min":30,"name":"tokenKey","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"exceptDomains":null,"hidden":false,"id":"_pbf_auth_email_","name":"email","onlyDomains":null,"presentable":false,"required":false,"system":true,"type":"email"},{"hidden":false,"id":"_pbf_auth_emailVisibility_","name":"emailVisibility","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"_pbf_auth_verified_","name":"verified","presentable":false,"required":false,"system":true,"type":"bool"},{"autogeneratePattern":"users[0-9]{5}","hidden":false,"id":"_pbf_auth_username_","max":150,"min":3,"name":"username","pattern":"^[\\\\w][\\\\w\\\\.\\\\-]*$","presentable":false,"primaryKey":false,"required":true,"system":false,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"users_name","max":0,"min":0,"name":"name","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":false,"type":"text"},{"hidden":false,"id":"users_avatar","maxSelect":1,"maxSize":5242880,"mimeTypes":["image/jpg","image/jpeg","image/png","image/svg+xml","image/gif"],"name":"avatar","presentable":false,"protected":false,"required":false,"system":false,"thumbs":["70x50","70x50t","70x50b","70x50f","0x50","70x0"],"type":"file"},{"hidden":false,"id":"xtecur3m","maxSelect":5,"maxSize":5242880,"mimeTypes":null,"name":"file","presentable":false,"protected":false,"required":false,"system":false,"thumbs":null,"type":"file"},{"cascadeDelete":false,"collectionId":"sz5l5z67tg7gku0","hidden":false,"id":"lkeigvv3","maxSelect":1,"minSelect":0,"name":"rel","presentable":false,"required":false,"system":false,"type":"relation"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]	\N	id = @request.auth.id		id = @request.auth.id	id = @request.auth.id	{"authRule":"","manageRule":null,"authAlert":{"enabled":true,"emailTemplate":{"subject":"Login from a new location","body":"\\u003cp\\u003eHello {RECORD:name}{RECORD:tokenKey},\\u003c/p\\u003e\\n\\u003cp\\u003eWe noticed a login to your {APP_NAME} account from a new location.\\u003c/p\\u003e\\n\\u003cp\\u003eIf this was you, you may disregard this email.\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003cstrong\\u003eIf this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}},"oauth2":{"providers":[{"pkce":null,"name":"gitlab","clientId":"test1","clientSecret":"test2","authURL":"","tokenURL":"","userInfoURL":"","displayName":""},{"pkce":null,"name":"google","clientId":"test","clientSecret":"test2","authURL":"","tokenURL":"","userInfoURL":"","displayName":""}],"mappedFields":{"id":"","name":"","username":"username","avatarURL":""},"enabled":true},"passwordAuth":{"enabled":true,"identityFields":["email","username"]},"mfa":{"enabled":true,"duration":1800,"rule":""},"otp":{"enabled":true,"duration":300,"length":8,"emailTemplate":{"subject":"OTP for {APP_NAME}","body":"\\u003cp\\u003eHello {RECORD:name}{RECORD:tokenKey},\\u003c/p\\u003e\\n\\u003cp\\u003eYour one-time password is: \\u003cstrong\\u003e{OTP}\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask for the one-time password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}},"authToken":{"secret":"PjVU4hAV7CZIWbCByJHkDcMUlSEWCLI6M5aWSZOpEq0a3rYxKT","duration":1209600},"passwordResetToken":{"secret":"BC6jYPe4JXpQGGNzu6VXtYw0yhKoH2mh2ezIJClOJQuZYrd4Ol","duration":1800},"emailChangeToken":{"secret":"eON2TTJZiGCEi7mvUvwMLADj8CMHQzwZN3gmyMjQb24EY08ATP","duration":1800},"verificationToken":{"secret":"dgGGHlzzdCJ2C5MjXGoondllwSXkJHyL50FuvLvXGHNmBhvGKO","duration":604800},"fileToken":{"secret":"4Ax9zDm2Rwtny81dGaGQrJQBnIx5wVOuNe89X6v7NbNzrAZhvn","duration":180},"verificationTemplate":{"subject":"Verify your {APP_NAME} email","body":"\\u003cp\\u003eHello {RECORD:name}{RECORD:tokenKey},\\u003c/p\\u003e\\n\\u003cp\\u003eThank you for joining us at {APP_NAME}.\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to verify your email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eVerify\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"resetPasswordTemplate":{"subject":"Reset your {APP_NAME} password","body":"\\u003cp\\u003eHello {RECORD:name}{RECORD:tokenKey},\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to reset your password.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eReset password\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to reset your password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"confirmEmailChangeTemplate":{"subject":"Confirm your {APP_NAME} new email address","body":"\\u003cp\\u003eHello {RECORD:name}{RECORD:tokenKey},\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to confirm your new email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eConfirm new email\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to change your email address, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}}	2022-10-10 09:49:46.145Z	2024-09-13 10:48:13.365Z	["CREATE UNIQUE INDEX \\"__pb_users_auth__username_idx\\" ON \\"users\\" (username)","CREATE UNIQUE INDEX \\"__pb_users_auth__email_idx\\" ON \\"users\\" (email) WHERE email != ''","CREATE UNIQUE INDEX \\"__pb_users_auth__tokenKey_idx\\" ON \\"users\\" (\\"tokenKey\\")","CREATE INDEX \\"__pb_users_auth__created_idx\\" ON \\"users\\" (\\"created\\")"]
wsmn24bux7wo113	f	base	demo1	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"u7spsiph","max":0,"min":0,"name":"text","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":false,"type":"text"},{"hidden":false,"id":"puk534el","name":"bool","presentable":false,"required":false,"system":false,"type":"bool"},{"exceptDomains":null,"hidden":false,"id":"ktas5n7b","name":"url","onlyDomains":null,"presentable":false,"required":false,"system":false,"type":"url"},{"hidden":false,"id":"dc4abz4i","maxSelect":1,"name":"select_one","presentable":false,"required":false,"system":false,"type":"select","values":["optionA","optionB","optionC"]},{"hidden":false,"id":"owtlq7zl","maxSelect":3,"name":"select_many","presentable":false,"required":false,"system":false,"type":"select","values":["optionA","optionB","optionC"]},{"hidden":false,"id":"4ulkdevf","maxSelect":1,"maxSize":5242880,"mimeTypes":null,"name":"file_one","presentable":false,"protected":true,"required":false,"system":false,"thumbs":null,"type":"file"},{"hidden":false,"id":"fjzhrsvq","maxSelect":99,"maxSize":5242880,"mimeTypes":null,"name":"file_many","presentable":false,"protected":false,"required":false,"system":false,"thumbs":null,"type":"file"},{"hidden":false,"id":"1z1ld0i5","max":null,"min":null,"name":"number","onlyInt":false,"presentable":false,"required":false,"system":false,"type":"number"},{"exceptDomains":null,"hidden":false,"id":"khvhpwgj","name":"email","onlyDomains":null,"presentable":false,"required":false,"system":false,"type":"email"},{"hidden":false,"id":"ro6p02gk","max":"","min":"","name":"datetime","presentable":false,"required":false,"system":false,"type":"date"},{"hidden":false,"id":"ei2fg4v1","maxSize":5242880,"name":"json","presentable":false,"required":false,"system":false,"type":"json"},{"cascadeDelete":false,"collectionId":"wsmn24bux7wo113","hidden":false,"id":"zaedritp","maxSelect":1,"minSelect":0,"name":"rel_one","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":true,"collectionId":"_pb_users_auth_","hidden":false,"id":"t9bpk2ug","maxSelect":9999,"minSelect":0,"name":"rel_many","presentable":false,"required":false,"system":false,"type":"relation"},{"hidden":false,"id":"geoPoint3081106212","name":"point","presentable":false,"required":false,"system":false,"type":"geoPoint"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]	\N	\N	\N	\N	\N	{}	2022-10-10 09:51:19.868Z	2025-04-02 11:04:29.260Z	["CREATE INDEX \\"_wsmn24bux7wo113_created_idx\\" ON \\"demo1\\" (\\"created\\")"]
sz5l5z67tg7gku0	f	base	demo2	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"mkrguaaf","max":0,"min":2,"name":"title","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":false,"type":"text"},{"hidden":false,"id":"izkl5z2s","name":"active","presentable":false,"required":false,"system":false,"type":"bool"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]						{}	2022-10-10 09:51:28.452Z	2023-03-20 16:39:25.211Z	["CREATE INDEX \\"idx_demo2_created\\" ON \\"demo2\\" (\\"created\\")","CREATE UNIQUE INDEX \\"idx_unique_demo2_title\\" on \\"demo2\\" (\\"title\\")","CREATE INDEX \\"idx_demo2_active\\" ON \\"demo2\\" (\\n\\"active\\"\\n)"]
wzlqyes4orhoygb	f	base	demo3	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"w5z2x0nq","max":0,"min":0,"name":"title","pattern":"","presentable":true,"primaryKey":false,"required":false,"system":false,"type":"text"},{"hidden":false,"id":"tgqrbwio","maxSelect":99,"maxSize":5242880,"mimeTypes":null,"name":"files","presentable":false,"protected":false,"required":false,"system":false,"thumbs":null,"type":"file"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]	@request.auth.id != "" && @request.auth.collectionName != "users"	@request.auth.id != "" && @request.auth.collectionName != "users"	@request.auth.id != "" && @request.auth.collectionName != "users"	@request.auth.id != "" && @request.auth.collectionName != "users"	@request.auth.id != "" && @request.auth.collectionName != "users"	{}	2022-10-10 09:51:36.853Z	2023-11-20 18:26:53.176Z	["CREATE INDEX \\"_wzlqyes4orhoygb_created_idx\\" ON \\"demo3\\" (\\"created\\")"]
4d1blo5cuycfaca	f	base	demo4	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"erkxnabw","max":0,"min":0,"name":"title","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":false,"type":"text"},{"cascadeDelete":false,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"t5jskeyz","maxSelect":1,"minSelect":0,"name":"rel_one_no_cascade","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"ldrcjhk8","maxSelect":1,"minSelect":0,"name":"rel_one_no_cascade_required","presentable":false,"required":true,"system":false,"type":"relation"},{"cascadeDelete":true,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"pl5lcd4y","maxSelect":1,"minSelect":0,"name":"rel_one_cascade","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"jz0oue3z","maxSelect":999,"minSelect":0,"name":"rel_many_no_cascade","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"bsxyqrhb","maxSelect":999,"minSelect":0,"name":"rel_many_no_cascade_required","presentable":false,"required":true,"system":false,"type":"relation"},{"cascadeDelete":true,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"kwmchnf7","maxSelect":999,"minSelect":0,"name":"rel_many_cascade","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"pmynkqk5","maxSelect":1,"minSelect":0,"name":"rel_one_unique","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"wzlqyes4orhoygb","hidden":false,"id":"mjzyk9vb","maxSelect":999,"minSelect":0,"name":"rel_many_unique","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"4d1blo5cuycfaca","hidden":false,"id":"dagiyxj4","maxSelect":1,"minSelect":0,"name":"self_rel_one","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"4d1blo5cuycfaca","hidden":false,"id":"tsrki8kc","maxSelect":999,"minSelect":0,"name":"self_rel_many","presentable":false,"required":false,"system":false,"type":"relation"},{"hidden":false,"id":"4wpx0hhx","maxSize":5242880,"name":"json_array","presentable":false,"required":false,"system":false,"type":"json"},{"hidden":false,"id":"ufpwiqnx","maxSize":5242880,"name":"json_object","presentable":false,"required":false,"system":false,"type":"json"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]			@request.auth.collectionName = 'users'	@request.auth.collectionName = 'users'		{}	2022-10-10 09:51:47.770Z	2024-09-01 12:15:28.332Z	["CREATE INDEX \\"_4d1blo5cuycfaca_created_idx\\" ON \\"demo4\\" (\\"created\\")","CREATE UNIQUE INDEX \\"idx_luoQV2A\\" ON \\"demo4\\" (\\"rel_one_unique\\") WHERE rel_one_unique != ''","CREATE UNIQUE INDEX \\"idx_IjL94ze\\" ON \\"demo4\\" ((\\"rel_many_unique\\" #>> '{}')) WHERE rel_many_unique::text != '[]'"]
v851q4r790rhknl	f	auth	clients	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"cost":10,"hidden":true,"id":"_pbf_auth_password_","max":0,"min":8,"name":"password","pattern":"","presentable":false,"required":true,"system":true,"type":"password"},{"autogeneratePattern":"[a-zA-Z0-9_]{50}","hidden":true,"id":"_pbf_auth_tokenKey_","max":60,"min":30,"name":"tokenKey","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"exceptDomains":null,"hidden":false,"id":"_pbf_auth_email_","name":"email","onlyDomains":null,"presentable":false,"required":true,"system":true,"type":"email"},{"hidden":false,"id":"_pbf_auth_emailVisibility_","name":"emailVisibility","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"_pbf_auth_verified_","name":"verified","presentable":false,"required":false,"system":true,"type":"bool"},{"autogeneratePattern":"users[0-9]{5}","hidden":false,"id":"_pbf_auth_username_","max":150,"min":3,"name":"username","pattern":"^[\\\\w][\\\\w\\\\.\\\\-]*$","presentable":false,"primaryKey":false,"required":true,"system":false,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"lacorw19","max":0,"min":0,"name":"name","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":false,"type":"text"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]	\N	\N	\N	\N	\N	{"authRule":"verified=true","manageRule":null,"oauth2":{"providers":[],"mappedFields":{"id":"","name":"","username":"username","avatarUrl":""},"enabled":false},"passwordAuth":{"enabled":true,"identityFields":["email","username"]},"mfa":{"enabled":false,"duration":1800},"otp":{"enabled":false,"duration":300,"length":8,"emailTemplate":{"subject":"OTP for {APP_NAME}","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eYour one-time password is: \\u003cstrong\\u003e{OTP}\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask for the one-time password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}},"authToken":{"secret":"PjVU4hAV7CZIWbCByJHkDcMUlSEWCLI6M5aWSZOpEq0a3rYxKT","duration":1209600},"passwordResetToken":{"secret":"BC6jYPe4JXpQGGNzu6VXtYw0yhKoH2mh2ezIJClOJQuZYrd4Ol","duration":1800},"emailChangeToken":{"secret":"eON2TTJZiGCEi7mvUvwMLADj8CMHQzwZN3gmyMjQb24EY08ATP","duration":1800},"verificationToken":{"secret":"dgGGHlzzdCJ2C5MjXGoondllwSXkJHyL50FuvLvXGHNmBhvGKO","duration":604800},"fileToken":{"secret":"4Ax9zDm2Rwtny81dGaGQrJQBnIx5wVOuNe89X6v7NbNzrAZhvn","duration":180},"verificationTemplate":{"subject":"Verify your {APP_NAME} email","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eThank you for joining us at {APP_NAME}.\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to verify your email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eVerify\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"resetPasswordTemplate":{"subject":"Reset your {APP_NAME} password","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to reset your password.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eReset password\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to reset your password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"confirmEmailChangeTemplate":{"subject":"Confirm your {APP_NAME} new email address","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to confirm your new email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eConfirm new email\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to change your email address, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"authAlert":{"enabled":true,"emailTemplate":{"subject":"Login from a new location","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eWe noticed a login to your {APP_NAME} account from a new location.\\u003c/p\\u003e\\n\\u003cp\\u003eIf this was you, you may disregard this email.\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003cstrong\\u003eIf this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}}}	2022-10-11 09:41:36.712Z	2024-06-28 15:59:08.446Z	["CREATE UNIQUE INDEX \\"_v851q4r790rhknl_username_idx\\" ON \\"clients\\" (username)","CREATE UNIQUE INDEX \\"_v851q4r790rhknl_email_idx\\" ON \\"clients\\" (email) WHERE email != ''","CREATE UNIQUE INDEX \\"_v851q4r790rhknl_tokenKey_idx\\" ON \\"clients\\" (\\"tokenKey\\")","CREATE INDEX \\"_v851q4r790rhknl_created_idx\\" ON \\"clients\\" (\\"created\\")"]
kpv709sk2lqbqk8	t	auth	nologin	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"cost":10,"hidden":true,"id":"_pbf_auth_password_","max":0,"min":8,"name":"password","pattern":"","presentable":false,"required":true,"system":true,"type":"password"},{"autogeneratePattern":"[a-zA-Z0-9_]{50}","hidden":true,"id":"_pbf_auth_tokenKey_","max":60,"min":30,"name":"tokenKey","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"exceptDomains":null,"hidden":false,"id":"_pbf_auth_email_","name":"email","onlyDomains":null,"presentable":false,"required":true,"system":true,"type":"email"},{"hidden":false,"id":"_pbf_auth_emailVisibility_","name":"emailVisibility","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"_pbf_auth_verified_","name":"verified","presentable":false,"required":false,"system":true,"type":"bool"},{"autogeneratePattern":"users[0-9]{5}","hidden":false,"id":"_pbf_auth_username_","max":150,"min":3,"name":"username","pattern":"^[\\\\w][\\\\w\\\\.\\\\-]*$","presentable":false,"primaryKey":false,"required":true,"system":false,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"x8zzktwe","max":0,"min":0,"name":"name","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":false,"type":"text"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]						{"authRule":"","manageRule":"@request.auth.collectionName = 'users'","oauth2":{"providers":[{"name":"gitlab","pkce":null,"clientId":"test","clientSecret":"test","authUrl":"","tokenUrl":"","userApiUrl":"","displayName":""}],"mappedFields":{"id":"","name":"","username":"username","avatarUrl":""},"enabled":false},"passwordAuth":{"enabled":false,"identityFields":["email"]},"mfa":{"enabled":false,"duration":1800},"otp":{"enabled":false,"duration":300,"length":8,"emailTemplate":{"subject":"OTP for {APP_NAME}","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eYour one-time password is: \\u003cstrong\\u003e{OTP}\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask for the one-time password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}},"authToken":{"secret":"PjVU4hAV7CZIWbCByJHkDcMUlSEWCLI6M5aWSZOpEq0a3rYxKT","duration":1209600},"passwordResetToken":{"secret":"BC6jYPe4JXpQGGNzu6VXtYw0yhKoH2mh2ezIJClOJQuZYrd4Ol","duration":1800},"emailChangeToken":{"secret":"eON2TTJZiGCEi7mvUvwMLADj8CMHQzwZN3gmyMjQb24EY08ATP","duration":1800},"verificationToken":{"secret":"dgGGHlzzdCJ2C5MjXGoondllwSXkJHyL50FuvLvXGHNmBhvGKO","duration":604800},"fileToken":{"secret":"4Ax9zDm2Rwtny81dGaGQrJQBnIx5wVOuNe89X6v7NbNzrAZhvn","duration":180},"verificationTemplate":{"subject":"Verify your {APP_NAME} email","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eThank you for joining us at {APP_NAME}.\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to verify your email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eVerify\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"resetPasswordTemplate":{"subject":"Reset your {APP_NAME} password","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to reset your password.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eReset password\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to reset your password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"confirmEmailChangeTemplate":{"subject":"Confirm your {APP_NAME} new email address","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to confirm your new email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eConfirm new email\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to change your email address, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"authAlert":{"enabled":true,"emailTemplate":{"subject":"Login from a new location","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eWe noticed a login to your {APP_NAME} account from a new location.\\u003c/p\\u003e\\n\\u003cp\\u003eIf this was you, you may disregard this email.\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003cstrong\\u003eIf this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}}}	2022-10-12 10:39:21.294Z	2024-06-26 18:37:35.353Z	["CREATE UNIQUE INDEX \\"_kpv709sk2lqbqk8_username_idx\\" ON \\"nologin\\" (username)","CREATE UNIQUE INDEX \\"_kpv709sk2lqbqk8_email_idx\\" ON \\"nologin\\" (email) WHERE email != ''","CREATE UNIQUE INDEX \\"_kpv709sk2lqbqk8_tokenKey_idx\\" ON \\"nologin\\" (\\"tokenKey\\")","CREATE INDEX \\"_kpv709sk2lqbqk8_created_idx\\" ON \\"nologin\\" (\\"created\\")"]
9n89pl5vkct6330	f	base	demo5	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"_pbf_text_id_","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"hidden":false,"id":"sozvpexq","maxSelect":1,"name":"select_one","presentable":false,"required":false,"system":false,"type":"select","values":["a","b","c","d"]},{"hidden":false,"id":"qlq1nxlc","maxSelect":5,"name":"select_many","presentable":false,"required":false,"system":false,"type":"select","values":["a","b","c","d","e"]},{"cascadeDelete":false,"collectionId":"4d1blo5cuycfaca","hidden":false,"id":"ajrrsq1a","maxSelect":1,"minSelect":0,"name":"rel_one","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":false,"collectionId":"4d1blo5cuycfaca","hidden":false,"id":"soxhs0ou","maxSelect":5,"minSelect":0,"name":"rel_many","presentable":false,"required":false,"system":false,"type":"relation"},{"hidden":false,"id":"kvbyzuqj","max":null,"min":null,"name":"total","onlyInt":false,"presentable":false,"required":false,"system":false,"type":"number"},{"hidden":false,"id":"ob7dsrcl","maxSelect":1,"maxSize":5242880,"mimeTypes":null,"name":"file","presentable":false,"protected":false,"required":false,"system":false,"thumbs":null,"type":"file"},{"hidden":false,"id":"_pbf_autodate_created_","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"},{"hidden":false,"id":"_pbf_autodate_updated_","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":false,"type":"autodate"}]	select_many:length = 3	rel_many.self_rel_many.rel_many_cascade.files:length = 1	@request.body.total = 3	@request.body.total = 3	@request.query.test:isset = true	{}	2023-01-07 13:13:08.733Z	2023-04-04 13:10:52.723Z	["CREATE INDEX \\"_9n89pl5vkct6330_created_idx\\" ON \\"demo5\\" (\\"created\\")"]
v9gwnfh02gjq1q0	f	view	view1	[{"autogeneratePattern":"","hidden":false,"id":"text3208210256","max":0,"min":0,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"_clone_F4wf","max":0,"min":0,"name":"text","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":false,"type":"text"},{"hidden":false,"id":"_clone_F2ir","name":"bool","presentable":false,"required":false,"system":false,"type":"bool"},{"exceptDomains":null,"hidden":false,"id":"_clone_X6zE","name":"url","onlyDomains":null,"presentable":false,"required":false,"system":false,"type":"url"},{"hidden":false,"id":"_clone_bZA2","maxSelect":1,"name":"select_one","presentable":false,"required":false,"system":false,"type":"select","values":["optionA","optionB","optionC"]},{"hidden":false,"id":"_clone_B5De","maxSelect":3,"name":"select_many","presentable":false,"required":false,"system":false,"type":"select","values":["optionA","optionB","optionC"]},{"hidden":false,"id":"_clone_gPsK","maxSelect":1,"maxSize":5242880,"mimeTypes":null,"name":"file_one","presentable":false,"protected":true,"required":false,"system":false,"thumbs":null,"type":"file"},{"hidden":false,"id":"_clone_sZpk","maxSelect":99,"maxSize":5242880,"mimeTypes":null,"name":"file_many","presentable":false,"protected":false,"required":false,"system":false,"thumbs":null,"type":"file"},{"hidden":false,"id":"_clone_Goev","max":null,"min":null,"name":"number","onlyInt":false,"presentable":false,"required":false,"system":false,"type":"number"},{"exceptDomains":null,"hidden":false,"id":"_clone_Das5","name":"email","onlyDomains":null,"presentable":false,"required":false,"system":false,"type":"email"},{"hidden":false,"id":"_clone_psUd","max":"","min":"","name":"datetime","presentable":false,"required":false,"system":false,"type":"date"},{"hidden":false,"id":"_clone_WfVC","maxSize":5242880,"name":"json","presentable":false,"required":false,"system":false,"type":"json"},{"cascadeDelete":false,"collectionId":"wsmn24bux7wo113","hidden":false,"id":"_clone_OqCt","maxSelect":1,"minSelect":0,"name":"rel_one","presentable":false,"required":false,"system":false,"type":"relation"},{"cascadeDelete":true,"collectionId":"_pb_users_auth_","hidden":false,"id":"_clone_9UR2","maxSelect":9999,"minSelect":0,"name":"rel_many","presentable":false,"required":false,"system":false,"type":"relation"},{"hidden":false,"id":"_clone_Y06D","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":false,"type":"autodate"}]	@request.auth.id != "" && bool = true	@request.auth.id != "" && bool = true	\N	\N	\N	{"viewQuery":"select id, text, bool, url, select_one, select_many, file_one, file_many, number, email, datetime, json, rel_one, rel_many, created from demo1"}	2023-02-12 18:58:12.315Z	2024-11-19 15:29:43.683Z	[]
ib3m2700k5hlsjz	f	view	view2	[{"autogeneratePattern":"","hidden":false,"id":"text3208210256","max":0,"min":0,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"hidden":false,"id":"_clone_77Ik","name":"state","presentable":false,"required":false,"system":false,"type":"bool"},{"hidden":false,"id":"_clone_IbSu","maxSelect":99,"maxSize":5242880,"mimeTypes":null,"name":"file_many","presentable":false,"protected":false,"required":false,"system":false,"thumbs":null,"type":"file"},{"cascadeDelete":true,"collectionId":"_pb_users_auth_","hidden":false,"id":"_clone_Asz0","maxSelect":9999,"minSelect":0,"name":"rel_many","presentable":false,"required":false,"system":false,"type":"relation"}]			\N	\N	\N	{"viewQuery":"SELECT view1.id, view1.bool as state, view1.file_many, view1.rel_many from view1\\n"}	2023-02-17 19:42:54.278Z	2024-11-19 15:29:43.857Z	[]
zahsr9d2mix2fvk	f	view	numeric_id_view	[{"autogeneratePattern":"","hidden":false,"id":"text3208210256","max":0,"min":0,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"exceptDomains":null,"hidden":false,"id":"_clone_jO9O","name":"email","onlyDomains":null,"presentable":false,"required":true,"system":true,"type":"email"}]			\N	\N	\N	{"viewQuery":"select (ROW_NUMBER() OVER()) as id, email from clients"}	2023-08-11 09:41:00.997Z	2024-11-19 15:29:43.901Z	[]
pbc_3142635823	t	auth	_superusers	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"text3208210256","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"cost":0,"hidden":true,"id":"password901924565","max":0,"min":8,"name":"password","pattern":"","presentable":false,"required":true,"system":true,"type":"password"},{"autogeneratePattern":"[a-zA-Z0-9_]{50}","hidden":true,"id":"text2504183744","max":60,"min":30,"name":"tokenKey","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"exceptDomains":null,"hidden":false,"id":"email3885137012","name":"email","onlyDomains":null,"presentable":false,"required":true,"system":true,"type":"email"},{"hidden":false,"id":"bool1547992806","name":"emailVisibility","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"bool256245529","name":"verified","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"autodate2990389176","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":true,"type":"autodate"},{"hidden":false,"id":"autodate3332085495","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":true,"type":"autodate"}]	\N	\N	\N	\N	\N	{"authRule":"","manageRule":null,"authAlert":{"enabled":true,"emailTemplate":{"subject":"Login from a new location","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eWe noticed a login to your {APP_NAME} account from a new location.\\u003c/p\\u003e\\n\\u003cp\\u003eIf this was you, you may disregard this email.\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003cstrong\\u003eIf this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}},"oauth2":{"providers":[],"mappedFields":{"id":"","name":"","username":"","avatarURL":""},"enabled":false},"passwordAuth":{"enabled":true,"identityFields":["email"]},"mfa":{"enabled":false,"duration":1800,"rule":""},"otp":{"enabled":false,"duration":300,"length":8,"emailTemplate":{"subject":"OTP for {APP_NAME}","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eYour one-time password is: \\u003cstrong\\u003e{OTP}\\u003c/strong\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask for the one-time password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}},"authToken":{"secret":"MyN3nDlzmHnuCjd35vb6cyIdqNr7Os0PmgiPVDMxmbFToSpBvS","duration":1209600},"passwordResetToken":{"secret":"fPSpFm9rxjj4mdeWYfyQ5OZQ4UWpyainTO0dqrJe3LHEYEDduq","duration":1800},"emailChangeToken":{"secret":"unYNiYeuIxH7BCV09NIb81abe2bkPgaexMYdDQ6uOOIFh74urD","duration":1800},"verificationToken":{"secret":"uhr68rXLVjPBWALFtw8uEHeQwDdN4t0MiTLr2pBWVkEQnNICe1","duration":259200},"fileToken":{"secret":"sjJAjTNPrOcRDmnIKwQm7qY9FyjuXTG5KNcaqw4U1TSDVfu4r9","duration":180},"verificationTemplate":{"subject":"Verify your {APP_NAME} email","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eThank you for joining us at {APP_NAME}.\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to verify your email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eVerify\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"resetPasswordTemplate":{"subject":"Reset your {APP_NAME} password","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to reset your password.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eReset password\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to reset your password, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"},"confirmEmailChangeTemplate":{"subject":"Confirm your {APP_NAME} new email address","body":"\\u003cp\\u003eHello,\\u003c/p\\u003e\\n\\u003cp\\u003eClick on the button below to confirm your new email address.\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  \\u003ca class=\\"btn\\" href=\\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\\" target=\\"_blank\\" rel=\\"noopener\\"\\u003eConfirm new email\\u003c/a\\u003e\\n\\u003c/p\\u003e\\n\\u003cp\\u003e\\u003ci\\u003eIf you didn't ask to change your email address, you can ignore this email.\\u003c/i\\u003e\\u003c/p\\u003e\\n\\u003cp\\u003e\\n  Thanks,\\u003cbr/\\u003e\\n  {APP_NAME} team\\n\\u003c/p\\u003e"}}	2024-06-20 09:29:22.826Z	2024-09-08 10:49:38.496Z	["CREATE UNIQUE INDEX \\"idx_tokenKey__pbc_3323866339\\" ON \\"_superusers\\" (\\"tokenKey\\")","CREATE UNIQUE INDEX \\"idx_email__pbc_3323866339\\" ON \\"_superusers\\" (\\"email\\") WHERE \\"email\\" != ''"]
pbc_2281828961	t	base	_externalAuths	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"text3208210256","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text455797646","max":0,"min":0,"name":"collectionRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text127846527","max":0,"min":0,"name":"recordRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text2462348188","max":0,"min":0,"name":"provider","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text1044722854","max":0,"min":0,"name":"providerId","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"hidden":false,"id":"autodate2990389176","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":true,"type":"autodate"},{"hidden":false,"id":"autodate3332085495","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":true,"type":"autodate"}]	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	\N	\N	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	{}	2024-06-20 09:29:22.862Z	2024-07-31 16:33:51.596Z	["CREATE UNIQUE INDEX \\"idx_externalAuths_record_provider\\" ON \\"externalAuths\\" (\\"collectionRef\\", \\"recordRef\\", provider)","CREATE UNIQUE INDEX \\"idx_externalAuths_collection_provider\\" ON \\"externalAuths\\" (\\"collectionRef\\", provider, \\"providerId\\")"]
pbc_2279338944	t	base	_mfas	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"text3208210256","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text455797646","max":0,"min":0,"name":"collectionRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text127846527","max":0,"min":0,"name":"recordRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text1582905952","max":0,"min":0,"name":"method","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"hidden":false,"id":"autodate2990389176","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":true,"type":"autodate"},{"hidden":false,"id":"autodate3332085495","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":true,"type":"autodate"}]	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	\N	\N	\N	{}	2024-06-20 09:29:22.894Z	2024-10-24 18:34:35.352Z	["CREATE INDEX \\"idx_mfas_collectionRef_recordRef\\" ON \\"mfas\\" (\\n  \\"collectionRef\\",\\n  \\"recordRef\\"\\n)"]
pbc_1638494021	t	base	_otps	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"text3208210256","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text455797646","max":0,"min":0,"name":"collectionRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text127846527","max":0,"min":0,"name":"recordRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"cost":8,"hidden":true,"id":"password901924565","max":0,"min":0,"name":"password","pattern":"","presentable":false,"required":true,"system":true,"type":"password"},{"hidden":false,"id":"autodate2990389176","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":true,"type":"autodate"},{"hidden":false,"id":"autodate3332085495","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":true,"type":"autodate"},{"autogeneratePattern":"","hidden":true,"id":"text3866985172","max":0,"min":0,"name":"sentTo","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":true,"type":"text"}]	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	\N	\N	\N	{}	2024-06-20 09:29:22.901Z	2024-11-19 15:29:43.633Z	["CREATE INDEX \\"idx_otps_collectionRef_recordRef\\" ON \\"otps\\" (\\n  \\"collectionRef\\",\\n  \\"recordRef\\"\\n)"]
pbc_4275539003	t	base	_authOrigins	[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"text3208210256","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text455797646","max":0,"min":0,"name":"collectionRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text127846527","max":0,"min":0,"name":"recordRef","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"text4228609354","max":0,"min":0,"name":"fingerprint","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"},{"hidden":false,"id":"autodate2990389176","name":"created","onCreate":true,"onUpdate":false,"presentable":false,"system":true,"type":"autodate"},{"hidden":false,"id":"autodate3332085495","name":"updated","onCreate":true,"onUpdate":true,"presentable":false,"system":true,"type":"autodate"}]	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	\N	\N	@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId	{}	2024-06-20 12:10:52.542Z	2024-07-31 16:32:24.722Z	["CREATE UNIQUE INDEX \\"idx_authOrigins_unique_pairs\\" ON \\"authDevices\\" (\\"collectionRef\\", \\"recordRef\\", fingerprint)"]
\.


--
-- Data for Name: "_externalAuths"; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public."_externalAuths" ("collectionRef", created, id, provider, "providerId", "recordRef", updated) FROM stdin;
_pb_users_auth_	2023-01-01 01:02:03.456Z	clmflokuq1xl341	google	test123	4q1xlclmfloku33	2022-01-01 01:01:01.123Z
_pb_users_auth_	2022-01-01 01:01:01.123Z	dlmflokuq1xl342	gitlab	test123	4q1xlclmfloku33	2022-01-01 01:01:01.123Z
v851q4r790rhknl	2024-07-29 11:56:28.439Z	f1z5b3843pzc964	google	test456	gk390qegs4y47wn	2024-07-29 11:56:35.151Z
_pb_users_auth_	2024-07-29 11:57:28.598Z	5eto7nmys833164	github	test456	bgs820n361vj1qd	2024-07-29 11:57:28.598Z
\.


--
-- Data for Name: _mfas; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._mfas ("collectionRef", created, id, method, "recordRef", updated) FROM stdin;
\.


--
-- Data for Name: _migrations; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._migrations (file, applied) FROM stdin;
1640988000_init.go	1665395386
1673167670_multi_match_migrate.go	1674485474882542
1677152688_rename_authentik_to_oidc.go	1678107894196945
1679943780_normalize_single_multiple_values.go	1680612240224966
1679943781_add_indexes_column.go	1680612240225611
1685164450_check_fk.go	1691746422498469
1689579878_renormalize_single_multiple_values.go	1691746422518206
1690319366_reset_null_values.go	1691746422524214
1690454337_transform_relations_to_views.go	1691746422528322
1691746861_created_view3.js	1691746861006404
1691746871_updated_view3.js	1691746871088040
1691746903_updated_view3.js	1691746903261135
1691747509_updated_numeric_id_view.js	1691747509618098
1691747531_updated_numeric_id_view.js	1691747531135925
1691747757_updated_numeric_id_view.js	1691747757587507
1691747912_resave_views.go	1691748137353973
1691747913_resave_views.go	1700504813174477
1692609521_copy_display_fields.go	1700504813176541
1701496825_allow_single_oauth2_provider_in_multiple_auth_collections.go	1701851393247905
1702134272_set_default_json_max_size.go	1702153002153382
1717233556_v0.23_migrate.go	1718885452550907
1719427055_updated_nologin.js	1719427055398268
1719427146_updated_users.js	1719427146443227
1719427201_updated_users.js	1719427201183755
1719427256_updated_users.js	1719427256670615
1719590348_updated_clients.js	1719590348495815
1720595195_updated_superusers.js	1720595195991995
1720862691_updated_numeric_id_view.js	1720862691263317
1720862697_updated_numeric_id_view.js	1720862697604141
1720862718_updated_view1.js	1720862718468637
1720862724_updated_view1.js	1720862724486655
1720862730_updated_view2.js	1720862730874909
1720862738_updated_view2.js	1720862738769086
1720862814_updated_numeric_id_view.js	1720862814980193
1720862822_updated_numeric_id_view.js	1720862822105793
1720863228_updated_numeric_id_view.js	1720863228586333
1720863270_updated_numeric_id_view.js	1720863270461984
1720863324_updated_numeric_id_view.js	1720863324872674
1720865157_updated_numeric_id_view.js	1720865157366879
1720865165_updated_view1.js	1720865165743836
1720865172_updated_view2.js	1720865172373325
1720865176_updated_view2.js	1720865176952277
1720865181_updated_view1.js	1720865181999132
1720865200_updated_view2.js	1720865200364699
1721987625_updated_authOrigins.js	1721987625054154
1721987631_updated_authOrigins.js	1721987631008720
1721987781_updated_authOrigins.js	1721987781298849
1721987788_updated_authOrigins.js	1721987788166593
1722443467_updated_otps.js	1722443467352142
1722443484_updated_mfas.js	1722443484760600
1722443494_updated_mfas.js	1722443494222859
1722443544_updated_authOrigins.js	1722443544810664
1722443631_updated_externalAuths.js	1722443631679710
1722443663_updated_superusers.js	1722443663936374
1724072848_updated_superusers.js	1724072848003576
1725192864_updated_view2.js	1725192864550539
1725192864_updated_view1.js	1725192864554724
1725192864_updated_demo1.js	1725192864567727
1725192928_updated_demo4.js	1725192928513955
1725792578_updated__superusers.js	1725792578533745
1726224305_updated_users.js	1726224305917687
1726224493_updated_users.js	1726224493523356
1640988000_aux_init.go	1729794875351351
1717233557_v0.23_migrate2.go	1729794875412697
1717233558_v0.23_migrate3.go	1730738463660106
1717233559_v0.23_migrate4.go	1732030183936141
1743591869_updated_demo1.js	1743591869355669
\.


--
-- Data for Name: _otps; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._otps ("collectionRef", created, id, password, "recordRef", updated, "sentTo") FROM stdin;
\.


--
-- Data for Name: _params; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._params (created, id, updated, value) FROM stdin;
\.


--
-- Data for Name: _superusers; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public._superusers (created, email, "emailVisibility", id, password, "tokenKey", updated, verified) FROM stdin;
2022-10-10 09:50:06.449Z	test@example.com	f	sywbhecnh46rhm0	$2a$13$6PP1fpShEsQHpZElw4.BN.vMl/2CMIKc2kpumJM//T4qV4RGUgb6i	O4rvW9FSUyTA3xUuQmXR3wHF2db9bHs19nBHeSgVTxerOsTAl4	2022-10-10 09:50:06.449Z	t
2022-10-10 11:22:27.982Z	test2@example.com	f	sbmbsdb40jyxf7h	$2a$13$1U32ttl.v1VtQkIH.a.z2.VVrU2IFwkW41TnY2OgyMcjADv34ynMK	cvg1nk1dKRFlazQH8nCKuFYwczdReQx6ZJimxXvei0uDyTkgEb	2022-10-10 11:22:31.096Z	t
2022-10-10 11:22:50.693Z	test3@example.com	f	9q2trqumvlyr3bd	$2a$13$gqXkfs0WjqTtUNRjRDxAHuLul1sxmA.elEfYuL0KT0ef5cMDt9Fjm	ezLvEu7DRFtUp9BI6nxtXCpgtp7qWaNQLdD6dDwjIVB0mA0uUr	2022-10-10 11:22:54.907Z	t
2024-07-26 12:17:43.787Z	test4@example.com	f	q911776rrfy658l	$2a$10$0CUWH81tPWUhxKqksZnO/eDfJiWTTG/xg7v5Yg8RhcLIGWY3oeCHe	B4LS_3ZiIlL_TQ97_W_u2_82__W1_l64037w_x211_F_BM_J4_	2024-07-26 12:17:43.787Z	t
\.


--
-- Data for Name: clients; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.clients (created, email, "emailVisibility", id, name, password, "tokenKey", updated, username, verified) FROM stdin;
2022-10-11 09:42:15.113Z	test@example.com	f	gk390qegs4y47wn		$2a$13$utsLtiNVf226R7hhXXyOM.kHTyTcVLRNqTIrT44Fw8JZbxMhIK0uu	rMb1gUpn27s53t66gOGscSHfYsa272cgOgn4nhTZIl4fIC8XP8	2022-10-11 09:42:15.113Z	clients57772	t
2022-10-11 09:42:25.984Z	test2@example.com	f	o1y0dd0spd786md	test_name	$2a$13$D/8j2Q7NzN5g/INiVn8qPOa2O3qkyZj7U82CXOyqzBFV0B9OdWhvC	RunSD73nFfH3sNScreizPcZYNkiPls2YjmFYPNo73cWKsDnZVm	2022-10-14 11:44:33.150Z	clients43362	f
\.


--
-- Data for Name: demo1; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.demo1 (created, id, updated, text, bool, url, select_one, file_one, file_many, number, email, datetime, "json", rel_one, select_many, rel_many, point) FROM stdin;
2022-10-14 10:13:12.397Z	84nmscqy84lsi1t	2023-01-04 20:11:20.732Z	test	t	https://example.copm	optionB	test_d61b33QdDU.txt	["test_QZFjKjXchk.txt","300_WlbFWSGmW9.png","logo_vcfJJG5TAh.svg","test_MaWC6mWyrP.txt","test_tC1Yc87DfC.txt"]	123456	test@example.com	2022-10-01 12:00:00.000Z	[\r\n  1,\r\n  2,\r\n  3\r\n]		["optionB","optionC"]	["oap640cot4yru2s"]	{"lat":0,"lon":0}
2022-10-14 10:14:04.685Z	al1h9ijdeojtsjy	2025-04-02 11:05:26.292Z	test2	f		optionB	300_Jsjq7RdBgA.png	[]	456	test2@example.com		null	84nmscqy84lsi1t	["optionB"]	["bgs820n361vj1qd","4q1xlclmfloku33","oap640cot4yru2s"]	{"lon":23.333157,"lat":42.654318}
2022-10-14 10:36:21.012Z	imy661ixudk5izi	2025-04-02 11:05:03.433Z	lorem ipsum	f				[]	0			null		[]	[]	{"lon":-74.006015,"lat":40.712728}
\.


--
-- Data for Name: demo2; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.demo2 (created, id, title, updated, active) FROM stdin;
2022-10-12 11:42:51.509Z	llvuca81nly1qls	test1	2022-10-12 11:42:51.509Z	f
2022-10-12 11:42:55.076Z	achvryl401bhse3	test2	2022-10-14 10:52:46.726Z	t
2022-10-12 11:42:58.215Z	0yxhwia2amd8gec	test3	2022-10-14 10:52:49.596Z	t
\.


--
-- Data for Name: demo3; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.demo3 (created, id, title, updated, files) FROM stdin;
2022-10-14 11:41:09.381Z	1tmknxy2868d869	test1	2022-10-14 11:41:09.381Z	[]
2022-10-14 11:41:12.548Z	lcl9d87w22ml6jy	test2	2022-10-14 14:08:29.095Z	["test_FLurQTgrY8.txt","300_UhLKX91HVb.png"]
2022-10-14 11:41:15.689Z	7nwo8tuiatetxdm	test3	2022-10-14 14:08:16.072Z	["test_JnXeKEwgwr.txt"]
2022-10-14 11:41:18.877Z	mk5fmymtx4wsprk	test4	2022-10-14 14:08:06.446Z	["300_JdfBOieXAW.png"]
\.


--
-- Data for Name: demo4; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.demo4 (created, id, title, updated, rel_one_no_cascade, rel_one_no_cascade_required, rel_one_cascade, self_rel_one, "json_array", "json_object", rel_one_unique, rel_many_no_cascade, rel_many_no_cascade_required, rel_many_cascade, rel_many_unique, self_rel_many) FROM stdin;
2022-10-14 14:15:43.771Z	qzaqccwrmva4o1n	test1	2022-10-20 19:23:59.427Z	mk5fmymtx4wsprk	lcl9d87w22ml6jy	7nwo8tuiatetxdm	i9naidtvr6qsgb4	[\r\n  1\r\n]	{"a": 123}		["1tmknxy2868d869"]	["lcl9d87w22ml6jy","1tmknxy2868d869"]	["mk5fmymtx4wsprk"]	[]	["i9naidtvr6qsgb4","qzaqccwrmva4o1n"]
2022-10-14 17:35:18.647Z	i9naidtvr6qsgb4	test2	2022-10-20 19:23:39.645Z		lcl9d87w22ml6jy		qzaqccwrmva4o1n	[1, 2, 3]	{"a": {"b": "test"}}		[]	["7nwo8tuiatetxdm"]	[]	[]	[]
\.


--
-- Data for Name: demo5; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.demo5 (created, id, rel_many, rel_one, select_many, select_one, updated, total, file) FROM stdin;
2023-01-07 13:13:24.249Z	qjeql998mtp1azp	["qzaqccwrmva4o1n","i9naidtvr6qsgb4"]	i9naidtvr6qsgb4	["b","c","a"]	b	2023-02-18 11:51:08.128Z	0	logo_vcf_jjg5_tah_9MtIHytOmZ.svg
2023-01-07 13:13:28.337Z	la4y2w4o98acwuj	[]		[]		2023-02-18 11:50:46.943Z	2	300_uh_lkx91_hvb_Da8K5pl069.png
\.


--
-- Data for Name: nologin; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.nologin (created, email, "emailVisibility", id, name, password, "tokenKey", updated, username, verified) FROM stdin;
2022-10-12 10:39:47.633Z	test@example.com	f	dc49k6jgejn40h3	test	$2a$13$395zl/w/nUmeo.vov5t6SO0/RsoE4suJnKLN2SzP41sJCe6PBAePq	6mi4JvnX8uIxS7JiO8LWl150Af8mnAGWkWRImHL2YB4XhOYf9c	2022-10-12 10:39:47.633Z	test_username	f
2022-10-12 10:39:59.058Z	test2@example.com	t	phhq3wr65cap535		$2a$13$ldkBOjQbXIXP3.xvJ8AXJep5I3kMzmXWu7wC75mP5RN4qK7mwQIK.	9Bsj2ogZ5b0Q3daAKZ5ZZrUuVb7lWHO0YDD4d5TbPCLfYE3COY	2022-10-14 12:03:49.335Z	viewers74618	f
2022-10-14 12:16:52.515Z	test3@example.com	f	oos036e9xvqeexy		$2a$13$oXNjS0xVi4aYHSaQ10bTOOMRvVFBY2t474S3UjM8.5BNnL5WB1Gce	uYafEN2cDgImHKMW4SfZ6tayrENAJVVWOTDxFXatOs2AO0KEMY	2022-10-14 12:16:52.515Z	nologin84738	t
\.


--
-- Data for Name: sqlite_stat1; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.sqlite_stat1 (tbl, idx, stat) FROM stdin;
authOrigins	idx_authOrigins_unique_pairs	7 4 3 1
authOrigins	sqlite_autoindex_authOrigins_1	7 1
_params	sqlite_autoindex__params_1	1 1
externalAuths	idx_externalAuths_collection_provider	4 2 1 1
externalAuths	idx_externalAuths_record_provider	4 2 2 1
externalAuths	sqlite_autoindex_externalAuths_1	4 1
superusers	sqlite_autoindex_superusers_1	4 1
demo5	_9n89pl5vkct6330_created_idx	2 1
demo5	sqlite_autoindex_demo5_1	2 1
users	__pb_users_auth__created_idx	3 1
users	__pb_users_auth__tokenKey_idx	3 1
users	__pb_users_auth__email_idx	3 1
users	__pb_users_auth__username_idx	3 1
users	sqlite_autoindex_users_1	3 1
nologin	_kpv709sk2lqbqk8_created_idx	3 1
nologin	_kpv709sk2lqbqk8_tokenKey_idx	3 1
nologin	_kpv709sk2lqbqk8_email_idx	3 1
nologin	_kpv709sk2lqbqk8_username_idx	3 1
nologin	sqlite_autoindex_nologin_1	3 1
demo2	idx_demo2_active	3 2
demo2	idx_unique_demo2_title	3 1
demo2	idx_demo2_created	3 1
demo2	sqlite_autoindex_demo2_1	3 1
_migrations	sqlite_autoindex__migrations_1	54 1
demo3	_wzlqyes4orhoygb_created_idx	4 1
demo3	sqlite_autoindex_demo3_1	4 1
clients	_v851q4r790rhknl_created_idx	2 1
clients	_v851q4r790rhknl_tokenKey_idx	2 1
clients	_v851q4r790rhknl_email_idx	2 1
clients	_v851q4r790rhknl_username_idx	2 1
clients	sqlite_autoindex_clients_1	2 1
demo4	idx_IjL94ze	0 0
demo4	idx_luoQV2A	0 0
demo4	_4d1blo5cuycfaca_created_idx	2 1
demo4	sqlite_autoindex_demo4_1	2 1
_superusers	idx_email__pbc_3323866339	4 1
_superusers	idx_tokenKey__pbc_3323866339	4 1
_superusers	sqlite_autoindex__superusers_1	4 1
_authOrigins	idx_authOrigins_unique_pairs	6 3 2 1
_authOrigins	sqlite_autoindex__authOrigins_1	6 1
_externalAuths	idx_externalAuths_collection_provider	4 2 1 1
_externalAuths	idx_externalAuths_record_provider	4 2 2 1
_externalAuths	sqlite_autoindex__externalAuths_1	4 1
_collections	idx__collections_type	16 6
_collections	sqlite_autoindex__collections_2	16 1
_collections	sqlite_autoindex__collections_1	16 1
demo1	_wsmn24bux7wo113_created_idx	3 1
demo1	sqlite_autoindex_demo1_1	3 1
\.


--
-- Data for Name: sqlite_stat4; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.sqlite_stat4 (tbl, idx, neq, nlt, ndlt, sample) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.users (avatar, created, email, "emailVisibility", id, name, password, "tokenKey", updated, username, verified, file, rel) FROM stdin;
300_1SEi6Q6U72.png	2022-10-10 10:36:26.879Z	test@example.com	f	4q1xlclmfloku33	test1	$2a$13$uULN32dHTkWQIALJP1N54.22HW/K9/qkXqcmsfz2hOA1wGeuUDbfG	tfYe7rCTX4D2KuWQY3pJjBifgsrMbecyXBatEPjrSfGEGS2jh6	2022-10-12 11:46:10.490Z	users75657	f	[]	llvuca81nly1qls
	2022-10-10 10:36:59.438Z	test2@example.com	f	oap640cot4yru2s		$2a$13$RC6/uXsHWM1ZV1v0cPJLRuWPXxyNINDmUDIHTq1x1dM.K.TBgWzFK	AQbE30CNb8Ncwr6Sg0sfvDJGJuepriTJN24EHZqO5DsEBTk1kA	2022-10-11 18:37:20.744Z	test2_username	t	["test_kfd2wYLxkz.txt"]	
	2022-10-10 10:37:33.119Z	test3@example.com	t	bgs820n361vj1qd	test3	$2a$13$qzF1J0ePG5.fvBrm0fVrtez5RRBjPOUoezvYRbTQGAfsT85d4XH2K	x6vHUi00LvM5bFeGIpwXN9xuol8k1BknfTmlySQ7YQWoTLKOa7	2022-10-12 11:46:02.462Z	users69238	t	[]	0yxhwia2amd8gec
\.


--
-- Name: "_authOrigins" sqlite_autoindex__authOrigins_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public."_authOrigins"
    ADD CONSTRAINT "sqlite_autoindex__authOrigins_1" PRIMARY KEY (id);


--
-- Name: _collections sqlite_autoindex__collections_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public._collections
    ADD CONSTRAINT sqlite_autoindex__collections_1 PRIMARY KEY (id);


--
-- Name: "_externalAuths" sqlite_autoindex__externalAuths_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public."_externalAuths"
    ADD CONSTRAINT "sqlite_autoindex__externalAuths_1" PRIMARY KEY (id);


--
-- Name: _mfas sqlite_autoindex__mfas_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public._mfas
    ADD CONSTRAINT sqlite_autoindex__mfas_1 PRIMARY KEY (id);


--
-- Name: _migrations sqlite_autoindex__migrations_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public._migrations
    ADD CONSTRAINT sqlite_autoindex__migrations_1 PRIMARY KEY (file);


--
-- Name: _otps sqlite_autoindex__otps_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public._otps
    ADD CONSTRAINT sqlite_autoindex__otps_1 PRIMARY KEY (id);


--
-- Name: _params sqlite_autoindex__params_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public._params
    ADD CONSTRAINT sqlite_autoindex__params_1 PRIMARY KEY (id);


--
-- Name: _superusers sqlite_autoindex__superusers_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public._superusers
    ADD CONSTRAINT sqlite_autoindex__superusers_1 PRIMARY KEY (id);


--
-- Name: clients sqlite_autoindex_clients_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT sqlite_autoindex_clients_1 PRIMARY KEY (id);


--
-- Name: demo1 sqlite_autoindex_demo1_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.demo1
    ADD CONSTRAINT sqlite_autoindex_demo1_1 PRIMARY KEY (id);


--
-- Name: demo2 sqlite_autoindex_demo2_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.demo2
    ADD CONSTRAINT sqlite_autoindex_demo2_1 PRIMARY KEY (id);


--
-- Name: demo3 sqlite_autoindex_demo3_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.demo3
    ADD CONSTRAINT sqlite_autoindex_demo3_1 PRIMARY KEY (id);


--
-- Name: demo4 sqlite_autoindex_demo4_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.demo4
    ADD CONSTRAINT sqlite_autoindex_demo4_1 PRIMARY KEY (id);


--
-- Name: demo5 sqlite_autoindex_demo5_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.demo5
    ADD CONSTRAINT sqlite_autoindex_demo5_1 PRIMARY KEY (id);


--
-- Name: nologin sqlite_autoindex_nologin_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.nologin
    ADD CONSTRAINT sqlite_autoindex_nologin_1 PRIMARY KEY (id);


--
-- Name: users sqlite_autoindex_users_1; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT sqlite_autoindex_users_1 PRIMARY KEY (id);


--
-- Name: _4d1blo5cuycfaca_created_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "_4d1blo5cuycfaca_created_idx" ON public.demo4 USING btree (created);


--
-- Name: _9n89pl5vkct6330_created_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "_9n89pl5vkct6330_created_idx" ON public.demo5 USING btree (created);


--
-- Name: __pb_users_auth__created_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "__pb_users_auth__created_idx" ON public.users USING btree (created);


--
-- Name: __pb_users_auth__email_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "__pb_users_auth__email_idx" ON public.users USING btree (email) WHERE "email" != '';


--
-- Name: __pb_users_auth__tokenKey_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "__pb_users_auth__tokenKey_idx" ON public.users USING btree ("tokenKey");


--
-- Name: __pb_users_auth__username_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "__pb_users_auth__username_idx" ON public.users USING btree (username);


--
-- Name: _kpv709sk2lqbqk8_created_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "_kpv709sk2lqbqk8_created_idx" ON public.nologin USING btree (created);


--
-- Name: _kpv709sk2lqbqk8_email_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "_kpv709sk2lqbqk8_email_idx" ON public.nologin USING btree (email);


--
-- Name: _kpv709sk2lqbqk8_tokenKey_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "_kpv709sk2lqbqk8_tokenKey_idx" ON public.nologin USING btree ("tokenKey");


--
-- Name: _kpv709sk2lqbqk8_username_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "_kpv709sk2lqbqk8_username_idx" ON public.nologin USING btree (username);


--
-- Name: _v851q4r790rhknl_created_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "_v851q4r790rhknl_created_idx" ON public.clients USING btree (created);


--
-- Name: _v851q4r790rhknl_email_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "_v851q4r790rhknl_email_idx" ON public.clients USING btree (email);


--
-- Name: _v851q4r790rhknl_tokenKey_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "_v851q4r790rhknl_tokenKey_idx" ON public.clients USING btree ("tokenKey");


--
-- Name: _v851q4r790rhknl_username_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "_v851q4r790rhknl_username_idx" ON public.clients USING btree (username);


--
-- Name: _wsmn24bux7wo113_created_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "_wsmn24bux7wo113_created_idx" ON public.demo1 USING btree (created);


--
-- Name: _wzlqyes4orhoygb_created_idx; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "_wzlqyes4orhoygb_created_idx" ON public.demo3 USING btree (created);


--
-- Name: idx__collections_type; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "idx__collections_type" ON public._collections USING btree (type);


--
-- Name: idx_authOrigins_unique_pairs; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "idx_authOrigins_unique_pairs" ON public."_authOrigins" USING btree ("collectionRef", "recordRef", fingerprint);


--
-- Name: idx_demo2_active; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "idx_demo2_active" ON public.demo2 USING btree (active);


--
-- Name: idx_demo2_created; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "idx_demo2_created" ON public.demo2 USING btree (created);


--
-- Name: idx_email__pbc_3323866339; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "idx_email__pbc_3323866339" ON public._superusers USING btree (email);


--
-- Name: idx_externalAuths_collection_provider; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "idx_externalAuths_collection_provider" ON public."_externalAuths" USING btree ("collectionRef", provider, "providerId");


--
-- Name: idx_externalAuths_record_provider; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "idx_externalAuths_record_provider" ON public."_externalAuths" USING btree ("collectionRef", "recordRef", provider);


--
-- Name: idx_mfas_collectionRef_recordRef; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "idx_mfas_collectionRef_recordRef" ON public._mfas USING btree ("collectionRef", "recordRef");


--
-- Name: idx_otps_collectionRef_recordRef; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX "idx_otps_collectionRef_recordRef" ON public._otps USING btree ("collectionRef", "recordRef");


--
-- Name: idx_tokenKey__pbc_3323866339; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "idx_tokenKey__pbc_3323866339" ON public._superusers USING btree ("tokenKey");


--
-- Name: idx_unique_demo2_title; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "idx_unique_demo2_title" ON public.demo2 USING btree (title);


--
-- Name: sqlite_autoindex__collections_2; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX "sqlite_autoindex__collections_2" ON public._collections USING btree (name);


--
-- PostgreSQL database dump complete
--



CREATE VIEW "view1" AS SELECT * FROM (select id, text, bool, url, select_one, select_many, file_one, file_many, number, email, datetime, json, rel_one, rel_many, created from demo1);
CREATE VIEW "view2" AS SELECT * FROM (SELECT view1.id, view1.bool as state, view1.file_many, view1.rel_many from view1);
CREATE VIEW "numeric_id_view" AS SELECT * FROM (SELECT CAST("id" as TEXT) "id","email" FROM (select (ROW_NUMBER() OVER()) as id, email from clients));