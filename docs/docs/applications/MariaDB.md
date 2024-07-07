# MariaDB
MariaDB is a free and open source database engine that started as a fork of MySQL and remains compatible with it. Many Linux distros install MariaDB by default if you install their `mysql` package. Podinate also recommends MariaDB over MySQL. 

## Installing
MySQL is available in the default Podinate package repositories. It can be installed directly.
```bash
podinate install mariadb
```
This will install mariadb with a randomly generated user and root password which you can retrieve from the pod's environment variables. 

More likely, you will want to use MariaDB as a sub-package when installing another application. Here is the most basic configuration:

!!! example 
    ```hcl
    package "database" {
        from = "mariadb"
        // Optional, will default to a randomly generated name
        //database_name = "my-cool-app"
    }

    pod "my-app" {
        name = "My Super App"
        image = "my-app"
        tag = "1.2.3"
        environment "MARIADB_USERNAME" {
            value = package.database.username
        }
        environment "MARIADB_PASSWORD" {
            value = package.database.password
        }
        environment "MARIADB_DATABASE" {
            value = package.database.database_name
        }
    }

    ```

## Performance
MariaDB is a high-performance and battle tested database engine, you should only run into performance issues if you are running a massive or complicated database. The first thing to try if you are running into performance issues should be to determine what indexes need to be added to your database schema. Beyond that, most apps tend to have one or two very common queries that can create excess load on the database. In these situations, consider checking out the [ReadySet](https://readyset.io/) database caching engine, which supports both MariaDB(MySQL) and Postgres. 

!!! note
    By default, ReadySet does not cache anything. For safety, you must log into its provided web interface and enable caching for specific queries. 

You can install ReadySet MariaDB with a single command. This will create a MariaDB database as before, but now with the ReadySet engine configured in front of the underlying MariaDB instance. 
```bash
podinate install readyset-mariadb
```





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