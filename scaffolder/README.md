# Scaffolder

##### CLI tool for scaffolding the Golang project

## Templates 
The Scaffold project uses Go Modules.
The whole scaffold project uses these modules structure:
- `github.com/go-chi/chi`
- `github.com/go-chi/cors`
- `github.com/go-ozzo/ozzo-validation`
- `github.com/lancer-kit/armory`
- `github.com/lancer-kit/uwe/v2`
- `github.com/rubenv/sql-migrate`
- `github.com/sirupsen/logrus`
- `github.com/urfave/cli`
- `gopkg.in/yaml.v2`

Templates are use [lancer-kit](https://github.com/lancer-kit) `armory` package for DB and log management and
`uwe/v2` package to provider worker pool functionality in the project.

#### The scaffolding system provides 3 versions of scaffolding:

1. Base \
    It contains: 
    * base CLI structure for serving the service;
    * info and version of the app;
    * empty config structure for workers;
    * chief structure for adding new uwe workers and root application structure.

2. API `--api flag` \
    It adds to the Base structure additional API logic:
    * api server worker
    * api configuration
    * handler, cors and http router empty structure

3. DB `--db flag` \
   It adds to the Base structure additional DB logic:
   * Init functions for dbschema migration
   * CLI for db migration
   * Query structure for database initialization and usage

4. Sime uwe `flag --base_uwe` \
   It adds and empty uwe worker to the projects structure

The scaffolder may generate the go mods with project name and tidy them on `--gomods` flag.

## Example
```
go get github.com/lancer-kit/forge
```



CLI Usage
```
scaffolder [global options] command [command options] [arguments...]
```

`Gen` cmd - generates the scaffold project contains following options that described in the table:

|Option                      | Required  | Description                                                       | Default value         |
|----------------------------|-----------|-------------------------------------------------------------------|-----------------------|
|`--gomods`                  | no        | Initializes the go modules with module name in scaffold project   | `false`               |
|`--output value, -o value`  | yes       | Specifies output dir to scaffold the project                      | `"./out"`             |
|`--domain value, -d value`  | no        | Specifies project scaffold domain                                 |          -            |
|`--name value, -n value`    | yes       | Specifies project scaffold name                                   | `"scaffold/project"`  |
|`--api`                     | no        | Specifies generation of optional API service logic                | `false`               |
|`--db`                      | no        | Specifies generation of optional DB service logic                 | `false`               |
|`--base_uwe`                | no        | Specifies generation of optional simple uwe worker logic          | `false`               |

Usage of CLI and its options described in [/example](https://github.com/lancer-kit/forge/master/scaffolder/example/Makefile)


