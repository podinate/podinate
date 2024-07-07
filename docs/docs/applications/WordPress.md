# WordPress
WordPress is the world's most widely used CMS and blogging platform. As such it's believed to power over 50 million websites worldwide. WordPress uses [MariaDB](../MariaDB) as its database engine. 


## Running on Podinate
WordPress is available as a Podinate package. To install it run:
```bash
podinate install wordpress
```

There are no required values for the package. 

## Speed / Resource Usage 
By default the image uses Apache in prefork mode. It would be a lot more resource efficient to switch to MPM, or to php-fpm within the image. 

WordPress has many caching plugins that can help by caching entire pre-rendered pages for anonymous visitors, which can greatly decrease resource usage for example if a blog post goes viral. 

## Environment Variables 

`WORDPRESS_DB_HOST` - the hostname of the database server you want WordPress to use. On Podinate this should be the one word name of your MariaDB instance's associated service. 

`WORDPRESS_DB_USER` - the username to use to connect to the database

`WORDPRESS_DB_PASSWORD` - the WordPress database user's password

`WORDPRESS_DB_NAME` - The name of the database wordpress should use.

## Volumes
WordPress is installed at `/var/www/html`. Within this directory is all of the PHP code files that make up the WordPress application. There is also a folder wp-content which contains all content generated on the fly. 

`/var/www/html/wp-content` - mount a volume here to contain all generated and uploaded content, for example themes and uploaded files. 

`/var/www/html/wp-content/uploads` - if you're building a custom WordPress image to deploy, you'll probably want to include all the needed files in other directories, and just set up this volume to hold uploads.