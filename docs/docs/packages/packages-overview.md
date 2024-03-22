# Podinate Packages
Podinate packages are written in terraform hcl and distributed as git repos. 

Podinate releases packages in repositories, which are just http servers containing a bunch of `tar.zst` files containing packages, and an `index.json` file listing them. The `index.json` file is also signed by a pgp key, so the authenticity of all files can be verified. 

## Structure
The basic structure of a package looks like this:

- Package.tar.zst
    - package.yaml
    - Some `.pod` files containing package definitions
    - scripts/ 

## Repositories
A Podinate package repository is just an http file server. There should be an `index.json` file which describes all the packages available in the repository, and an `index.json.sig` file which contains a PGP signature of the index.json file. The format of the index.json file is as follows:
```json
{
    "packages": [
        {
            "name": "wordpress",
            "description": "The world's number 1 blog and CMS engine, powering [...]",
            "documentation": "https://docs.podinate.com/packages/available-packages/WordPress/",
            "versions": [
                {
                    "ids": ["5"],
                    "package": "wordpress/5.tar.zst",
                    "sha256": "ABC12345"
                },
                {
                    "ids": ["latest", "6"],
                    "package": "wordpress/latest.tar.zst",
                    "sha256": "CCCC123456"
                },
                
            ]
        }
    ]

}
```

## Generating Repo
The Podinate CLI can generate a repo from a bunch of folders containing 