-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler version: 1.0.5
-- PostgreSQL version: 15.0
-- Project Site: pgmodeler.io
-- Model Author: ---

-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: podinate | type: DATABASE --
-- DROP DATABASE IF EXISTS podinate;
CREATE DATABASE podinate;
-- ddl-end --


-- object: public.project | type: TABLE --
-- DROP TABLE IF EXISTS public.project CASCADE;
CREATE TABLE public.project (
	uuid uuid NOT NULL,
	id text NOT NULL,
	name text,
	account_uuid uuid,
	created timestamp DEFAULT current_timestamp,
	CONSTRAINT project_pk PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN public.project.id IS E'Unique identifier for the project within the user''s account';
-- ddl-end --
COMMENT ON COLUMN public.project.name IS E'Human readable / display name for the project';
-- ddl-end --
ALTER TABLE public.project OWNER TO postgres;
-- ddl-end --

-- object: public.account | type: TABLE --
-- DROP TABLE IF EXISTS public.account CASCADE;
CREATE TABLE public.account (
	uuid uuid NOT NULL,
	id text,
	name text,
	owner_uuid uuid,
	flags jsonb,
	created timestamp NOT NULL DEFAULT current_timestamp,
	CONSTRAINT unique_account_slug UNIQUE (id),
	CONSTRAINT account_pk PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN public.account.id IS E'The unique identifier for the account within the system.';
-- ddl-end --
COMMENT ON COLUMN public.account.name IS E'The human readable / display name of the account';
-- ddl-end --
ALTER TABLE public.account OWNER TO postgres;
-- ddl-end --

-- object: public.project_pods | type: TABLE --
-- DROP TABLE IF EXISTS public.project_pods CASCADE;
CREATE TABLE public.project_pods (
	uuid uuid NOT NULL,
	id text,
	name text,
	image text,
	tag text,
	project_uuid uuid,
	CONSTRAINT project_pods_pk PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN public.project_pods.id IS E'The unique name for the deployment in kubernetes, used as the kuberenetes name.';
-- ddl-end --
COMMENT ON COLUMN public.project_pods.name IS E'Human readable / display name for the pod';
-- ddl-end --
COMMENT ON COLUMN public.project_pods.image IS E'The OCI image for the pod to run';
-- ddl-end --
COMMENT ON COLUMN public.project_pods.tag IS E'The image tag to run';
-- ddl-end --
ALTER TABLE public.project_pods OWNER TO postgres;
-- ddl-end --

-- object: account_fk | type: CONSTRAINT --
-- ALTER TABLE public.project DROP CONSTRAINT IF EXISTS account_fk CASCADE;
ALTER TABLE public.project ADD CONSTRAINT account_fk FOREIGN KEY (account_uuid)
REFERENCES public.account (uuid) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: project_fk | type: CONSTRAINT --
-- ALTER TABLE public.project_pods DROP CONSTRAINT IF EXISTS project_fk CASCADE;
ALTER TABLE public.project_pods ADD CONSTRAINT project_fk FOREIGN KEY (project_uuid)
REFERENCES public.project (uuid) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: unique_project_slug_per_account | type: CONSTRAINT --
-- ALTER TABLE public.project DROP CONSTRAINT IF EXISTS unique_project_slug_per_account CASCADE;
ALTER TABLE public.project ADD CONSTRAINT unique_project_slug_per_account UNIQUE (account_uuid,id);
-- ddl-end --

-- object: public."user" | type: TABLE --
-- DROP TABLE IF EXISTS public."user" CASCADE;
CREATE TABLE public."user" (
	uuid uuid NOT NULL,
	id text NOT NULL,
	display_name text,
	created timestamp DEFAULT CURRENT_TIMESTAMP,
	flags jsonb,
	CONSTRAINT user_pk PRIMARY KEY (uuid)
);
-- ddl-end --
ALTER TABLE public."user" OWNER TO postgres;
-- ddl-end --

-- object: public.oauth_login | type: TABLE --
-- DROP TABLE IF EXISTS public.oauth_login CASCADE;
CREATE TABLE public.oauth_login (
	provider text NOT NULL,
	provider_id text NOT NULL,
	provider_username text,
	access_token text NOT NULL,
	refresh_token text NOT NULL,
	authorised_user uuid NOT NULL,
	CONSTRAINT oauth_login_pk PRIMARY KEY (provider,provider_id)
);
-- ddl-end --
ALTER TABLE public.oauth_login OWNER TO postgres;
-- ddl-end --

-- object: public.api_key | type: TABLE --
-- DROP TABLE IF EXISTS public.api_key CASCADE;
CREATE TABLE public.api_key (
	key text NOT NULL,
	name text NOT NULL,
	user_uuid uuid NOT NULL,
	issued timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires timestamp,
	last_used timestamp,
	decription text

);
-- ddl-end --
COMMENT ON COLUMN public.api_key.key IS E'The one-way hashed API key';
-- ddl-end --
COMMENT ON COLUMN public.api_key.name IS E'User-provided name';
-- ddl-end --
COMMENT ON COLUMN public.api_key.decription IS E'User provided description';
-- ddl-end --
ALTER TABLE public.api_key OWNER TO postgres;
-- ddl-end --

-- object: public.login_session | type: TABLE --
-- DROP TABLE IF EXISTS public.login_session CASCADE;
CREATE TABLE public.login_session (
	session_id uuid NOT NULL,
	key text NOT NULL,
	value bytea,
	CONSTRAINT composite_primary PRIMARY KEY (session_id,key)
);
-- ddl-end --
ALTER TABLE public.login_session OWNER TO postgres;
-- ddl-end --

-- object: owner_uuid | type: CONSTRAINT --
-- ALTER TABLE public.account DROP CONSTRAINT IF EXISTS owner_uuid CASCADE;
ALTER TABLE public.account ADD CONSTRAINT owner_uuid FOREIGN KEY (owner_uuid)
REFERENCES public."user" (uuid) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: authorised_user_uuid | type: CONSTRAINT --
-- ALTER TABLE public.oauth_login DROP CONSTRAINT IF EXISTS authorised_user_uuid CASCADE;
ALTER TABLE public.oauth_login ADD CONSTRAINT authorised_user_uuid FOREIGN KEY (authorised_user)
REFERENCES public."user" (uuid) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: api_key_user_uuid | type: CONSTRAINT --
-- ALTER TABLE public.api_key DROP CONSTRAINT IF EXISTS api_key_user_uuid CASCADE;
ALTER TABLE public.api_key ADD CONSTRAINT api_key_user_uuid FOREIGN KEY (user_uuid)
REFERENCES public."user" (uuid) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --


