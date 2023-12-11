# MariaDB
MariaDB is a free and open source database engine that started as a fork of MySQL and remains largely compatible with it. Many Linux distros install MariaDB by default if you install their `mysql` package. 

## Official Image
The official MariaDB image is `docker.io/mariadb`, but it is also packaged by a number of other maintainers such as Bitnami. This documentation focuses on the official image. 

## Environment Variables
To get the official image to work, you will need the following environment variables: 

`MARIADB_DATABASE` - Set the database to be created

`MARIADB_USER` - Set the username that will have full permissions to the created database

`MARIADB_PASSWORD` - Set the password for that user

And one of the following:

`MARIADB_RANDOM_ROOT_PASSWORD` - Recommended - Set a random root password and print it in the logs

`MARIADB_ROOT_PASSWORD` - Specify a root password

## Volumes
MariaDB stores all database data at `/var/lib/mysql`, so a volume should be mounted here to make it persistent. 