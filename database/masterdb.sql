-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler version: 1.0.6
-- PostgreSQL version: 16.0
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
	uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	main_provider text NOT NULL,
	id text NOT NULL,
	display_name text,
	avatar_url text,
	created timestamp DEFAULT CURRENT_TIMESTAMP,
	email text NOT NULL,
	flags jsonb,
	CONSTRAINT user_pk PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN public."user".main_provider IS E'The provider string for this username, eg github, gitlab, podinate';
-- ddl-end --
COMMENT ON COLUMN public."user".display_name IS E'The user''s human name';
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
	value text,
	CONSTRAINT composite_primary PRIMARY KEY (session_id,key)
);
-- ddl-end --
ALTER TABLE public.login_session OWNER TO postgres;
-- ddl-end --

-- object: public.policy_attachment | type: TABLE --
-- DROP TABLE IF EXISTS public.policy_attachment CASCADE;
CREATE TABLE public.policy_attachment (
	account_uuid uuid NOT NULL,
	resource_id text NOT NULL,
	policy_uuid uuid NOT NULL,
	date_attached timestamp DEFAULT CURRENT_TIMESTAMP,
	valid_until timestamp,
	attached_by uuid,
	CONSTRAINT compound_primary PRIMARY KEY (account_uuid,resource_id,policy_uuid)
);
-- ddl-end --
COMMENT ON COLUMN public.policy_attachment.attached_by IS E'The user uuid that originally attached this policy';
-- ddl-end --
ALTER TABLE public.policy_attachment OWNER TO postgres;
-- ddl-end --

-- object: public.policy | type: TABLE --
-- DROP TABLE IF EXISTS public.policy CASCADE;
CREATE TABLE public.policy (
	uuid uuid NOT NULL,
	account_uuid uuid,
	id text,
	current_revision smallint,
	content text,
	date_added timestamp DEFAULT CURRENT_TIMESTAMP,
	added_by uuid,
	notes text,
	CONSTRAINT uuid_primary_key PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN public.policy.account_uuid IS E'The account to which this policy belongs';
-- ddl-end --
COMMENT ON COLUMN public.policy.content IS E'The content of the active revision of the policy';
-- ddl-end --
COMMENT ON COLUMN public.policy.added_by IS E'UUID of the user who added this';
-- ddl-end --
COMMENT ON COLUMN public.policy.notes IS E'Space for user to write some notes about this policy';
-- ddl-end --
ALTER TABLE public.policy OWNER TO postgres;
-- ddl-end --

-- object: public.policy_version | type: TABLE --
-- DROP TABLE IF EXISTS public.policy_version CASCADE;
CREATE TABLE public.policy_version (
	uuid uuid NOT NULL,
	policy_uuid uuid,
	version_number smallint NOT NULL,
	content text,
	comment text NOT NULL,
	user_uuid uuid,
	date_made timestamp DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT uuid PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN public.policy_version.content IS E'The policy document itself';
-- ddl-end --
COMMENT ON COLUMN public.policy_version.comment IS E'Commit message for the revision';
-- ddl-end --
COMMENT ON COLUMN public.policy_version.user_uuid IS E'User who made the revision';
-- ddl-end --
ALTER TABLE public.policy_version OWNER TO postgres;
-- ddl-end --

-- object: owner_uuid | type: CONSTRAINT --
-- ALTER TABLE public.account DROP CONSTRAINT IF EXISTS owner_uuid CASCADE;
ALTER TABLE public.account ADD CONSTRAINT owner_uuid FOREIGN KEY (owner_uuid)
REFERENCES public."user" (uuid) MATCH SIMPLE
ON DELETE CASCADE ON UPDATE CASCADE;
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

-- object: account_uuid | type: CONSTRAINT --
-- ALTER TABLE public.policy_attachment DROP CONSTRAINT IF EXISTS account_uuid CASCADE;
ALTER TABLE public.policy_attachment ADD CONSTRAINT account_uuid FOREIGN KEY (account_uuid)
REFERENCES public.account (uuid) MATCH SIMPLE
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: policy_uuid | type: CONSTRAINT --
-- ALTER TABLE public.policy_attachment DROP CONSTRAINT IF EXISTS policy_uuid CASCADE;
ALTER TABLE public.policy_attachment ADD CONSTRAINT policy_uuid FOREIGN KEY (policy_uuid)
REFERENCES public.policy (uuid) MATCH SIMPLE
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: policy_uuid | type: CONSTRAINT --
-- ALTER TABLE public.policy_version DROP CONSTRAINT IF EXISTS policy_uuid CASCADE;
ALTER TABLE public.policy_version ADD CONSTRAINT policy_uuid FOREIGN KEY (policy_uuid)
REFERENCES public.policy (uuid) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --


