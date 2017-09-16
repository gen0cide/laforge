# laforge

Competition infrastructure management for the cloud.

![laforge](http://vignette4.wikia.nocookie.net/p__/images/0/04/GeordiLaForge.jpg/revision/latest?cb=20151015040039&path-prefix=protagonist)

## Description

LaForge is a command line tool I developed to automate creating competition infrastructure. It automates and abstracts much of the terraform config writing into CLI and YAML.

At this time, LaForge only works on OS X and Linux systems.

## Usage

If you are starting from scratch, download the binary and place it in your path.

After that, you can run:

```
$ laforge init
```

This will walk you through a wizard to create a LaForge base directory and save it to `~/.lf_home`. This file will be referenced by LaForge to find your base directory for all subsequent commands.

After that, you can list the environments:

```
$ laforge env ls
```

You should see no environments. You should create one now:

```
$ laforge env create
```

Afterwards, you can use the following commands to list and create a network and a host.

```
$ laforge network ls
$ laforge network create
$ laforge host ls
$ laforge host create
```

**NOTE** - Make sure to read through ALL the YAML files created by laforge. There are required settings in some of them.

When you are ready to deploy your configuration, start by issuing the `build` subcommand:

```
$ laforge build
```

This will generate the `terraform/infra.tf` file as well as the `terraform/scripts` files needed by terraform. If this step fails, you have a bug in your configurations, or you have discovered a laforge bug!

When you want to use `terraform`, laforge aliases some subcommands to help you know what to run. None of these commands will be executed, it just prints the command you should run to the console.

```
$ laforge tf plan
$ laforge tf apply
$ laforge tf destroy
$ laforge tf nuke
```

You might need to install terraform plugins (AWS provider for example). This is a one time thing and can be done by issuing a `terraform get` command in your `terraform` root.

## Install

Download a release binary to a place in your `$PATH`.

## Known Issues
There are a few unimplemented subcommands:

 * `app` - Will be for package management.
 * `ssh` - Will eventually allow direct connection to hosts.
 * `doctor` - Will eventually troubleshoot issues in your environment.
 * `cd` - Will eventually drop you into locations in your filesystem.
 * `tf`* - Partially completed. Eventually will wrap terraform for you.

## Contribution

1. Fork ([https://github.com/gen0cide/laforge/fork](https://github.com/gen0cide/laforge/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[gen0cide](https://github.com/gen0cide)
