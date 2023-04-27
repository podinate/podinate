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
	id uuid NOT NULL,
	slug smallint,
	name text,
	account_id uuid,
	CONSTRAINT project_pk PRIMARY KEY (id)
);
-- ddl-end --
COMMENT ON COLUMN project.slug IS E'The unique identifier for the project within the account, for example "corp-blog"';
-- ddl-end --
COMMENT ON COLUMN project.name IS E'Human-readable / display name for the project eg "My Super Cool Blog"';
-- ddl-end --

-- ddl-end --

-- object: account | type: TABLE --
-- DROP TABLE IF EXISTS account CASCADE;
CREATE TABLE account (
	id uuid NOT NULL,
	slug text,
	name text,
	CONSTRAINT account_pk PRIMARY KEY (id)
);
-- ddl-end --
COMMENT ON COLUMN account.slug IS E'The slug is the unique identifier for the account, for example "my-prod-account"';
-- ddl-end --
COMMENT ON COLUMN account.name IS E'The friendly / display name of the account';
-- ddl-end --

-- ddl-end --

-- object: project_pods | type: TABLE --
-- DROP TABLE IF EXISTS project_pods CASCADE;
CREATE TABLE project_pods (
	id uuid,
	slug text,
	name text,
	image text,
	tag text,
	project_id uuid

);
-- ddl-end --
COMMENT ON TABLE project_pods IS E'Holds the pods for each project - under the hood a pod = a kubernetes deployment';
-- ddl-end --
COMMENT ON COLUMN project_pods.slug IS E'Used as the kubernetes pod name eg "wordpress-backend", the human readable name will be separate';
-- ddl-end --
COMMENT ON COLUMN project_pods.name IS E'Human readable name for the pod, for example "Wordpress Backend"';
-- ddl-end --
COMMENT ON COLUMN project_pods.image IS E'The docker image to run for the pod eg "registry.podinate.com/myaccount/cool-project" or "wordpress"';
-- ddl-end --
COMMENT ON COLUMN project_pods.tag IS E'The image tag to run eg "latest" / "5.0"';
-- ddl-end --

-- ddl-end --

-- object: account_fk | type: CONSTRAINT --
-- ALTER TABLE project DROP CONSTRAINT IF EXISTS account_fk CASCADE;
ALTER TABLE project ADD CONSTRAINT account_fk FOREIGN KEY (account_id)
REFERENCES account (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: project_fk | type: CONSTRAINT --
-- ALTER TABLE project_pods DROP CONSTRAINT IF EXISTS project_fk CASCADE;
ALTER TABLE project_pods ADD CONSTRAINT project_fk FOREIGN KEY (project_id)
REFERENCES project (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --


