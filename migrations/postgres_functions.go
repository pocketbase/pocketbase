package migrations

import "github.com/pocketbase/dbx"

func createSQLiteEquivalentFunctions(db dbx.Builder) error {
	//PostgreSQL:
	// 1. Check existance
	sql := `SELECT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'uuid_generate_v7');`
	var exists bool
	if err := db.NewQuery(sql).Row(&exists); err != nil {
		return err
	} else if exists {
		// The function already exists, no need to create it again
		return nil
	}

	// Postgres:
	// 2. Create function
	funcDef := `
	-- Enable built-in pgcrypto extension to use gen_random_bytes function
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	-- Adding "nocase" collation to be compatible with SQLite's built-in "nocase" collation
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

	-- Create a json_query_or_null function that handles any types.
	CREATE OR REPLACE FUNCTION json_query_or_null(p_input jsonb, p_query text) RETURNS jsonb AS $$
		SELECT JSON_QUERY(p_input, p_query)
	$$ LANGUAGE sql IMMUTABLE;

	-- Create a json_query_or_null function that handles any types.
	CREATE OR REPLACE FUNCTION json_query_or_null(p_input anyelement, p_query text) RETURNS jsonb AS $$
	BEGIN
		RETURN JSON_QUERY(p_input::text::jsonb, p_query);
	EXCEPTION WHEN others THEN
		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql STABLE;
	`
	_, err := db.NewQuery(funcDef).Execute()
	return err
}
