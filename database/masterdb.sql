-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler version: 1.0.2
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

use podinate;

-- object: project | type: TABLE --
-- DROP TABLE IF EXISTS project CASCADE;
CREATE TABLE project (
	uuid uuid NOT NULL,
	id text,
	name text,
	account_uuid uuid,
	CONSTRAINT project_pk PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN project.id IS E'Unique string identifier for the project within the user''s account';
-- ddl-end --
COMMENT ON COLUMN project.name IS E'Human readable / display name for the project';
-- ddl-end --
ALTER TABLE project OWNER TO postgres;
-- ddl-end --

-- object: account | type: TABLE --
-- DROP TABLE IF EXISTS account CASCADE;
CREATE TABLE account (
	uuid uuid NOT NULL,
	id text,
	name text,
	CONSTRAINT unique_account_id UNIQUE (id),
	CONSTRAINT account_pk PRIMARY KEY (uuid)
);
-- ddl-end --
COMMENT ON COLUMN account.id IS E'The unique string identifier for the account within the system.';
-- ddl-end --
COMMENT ON COLUMN account.name IS E'The human readable / display name of the account';
-- ddl-end --
ALTER TABLE account OWNER TO postgres;
-- ddl-end --

-- object: project_pods | type: TABLE --
-- DROP TABLE IF EXISTS project_pods CASCADE;
CREATE TABLE project_pods (
	uuid uuid NOT NULL,
	id text,
	name text,
	image text,
	tag text,
	count int,
	ram int,
	project_uuid uuid,
	CONSTRAINT project_pods_pk PRIMARY KEY (id)
);
-- ddl-end --
COMMENT ON COLUMN project_pods.id IS E'The string identifier for the deployment in kubernetes, used as the kuberenetes name.';
-- ddl-end --
COMMENT ON COLUMN project_pods.name IS E'Human readable / display name for the pod';
-- ddl-end --
COMMENT ON COLUMN project_pods.image IS E'The OCI image for the pod to run';
-- ddl-end --
COMMENT ON COLUMN project_pods.tag IS E'The image tag to run';
-- ddl-end --
ALTER TABLE project_pods OWNER TO postgres;
-- ddl-end --

-- object: account_fk | type: CONSTRAINT --
-- ALTER TABLE project DROP CONSTRAINT IF EXISTS account_fk CASCADE;
ALTER TABLE project ADD CONSTRAINT account_fk FOREIGN KEY (account_uuid)
REFERENCES account (uuid) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: project_fk | type: CONSTRAINT --
-- ALTER TABLE project_pods DROP CONSTRAINT IF EXISTS project_fk CASCADE;
ALTER TABLE project_pods ADD CONSTRAINT project_fk FOREIGN KEY (project_uuid)
REFERENCES project (uuid) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: unique_project_slug_per_account | type: CONSTRAINT --
-- ALTER TABLE project DROP CONSTRAINT IF EXISTS unique_project_slug_per_account CASCADE;
ALTER TABLE project ADD CONSTRAINT unique_project_id_per_account UNIQUE (account_uuid,id);
-- ddl-end --


