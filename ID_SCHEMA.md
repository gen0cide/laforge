# Laforge ID & Schemas

The ID can be 4-32 lowercase alpha-numeric characters and can include `-` (hypens) as long as it is not the first or last character.
IDs must not match any string in the ID blacklist.
This conforms to include MD5, UUID, xxHash, and other common ID formats.

This format is then standardizes an access model for key access within the standard in a similar way to common folder
structures. This will "chroot" the root `/` as the anchor in which base is located. So all subsequent unique keys will
corraspond to their relative path within the base context.

Example:

```
/ <~ root, the competition base where base.laforge is located.
/envs/example-env <~ example path key using the environment's compliant ID.
```

Global types are isolated into their own namespace and maybe have subdirectories. Those can be represented in their path. As an example, lets say you had the following files relative to your base:

```
/scripts
/scripts/linux/ubuntu/install-mysql.sh
/scripts/linux/ubuntu/install-mysql.sh.laforge
```

Inside the laforge file referenced above, you'd see the following:

```
script "/scripts/linux/ubuntu/install-mysql" {
  source = "./install-mysql.sh"
  // ...
}
```

Note the file extension has been dropped from the ID definition and represents a fully qualified path from the BaseContext's root directory. This subsequently can be referenced any place you need a fully qualified script ID. For example, inside of a host configuration, you could reference the following:

```
  provisioning_steps = [
    "/scripts/linux/ubuntu/install-mysql"
  ]
```

This standardization will allow for much more rich features in the future such as adding wildcard and pattern matching identification into fields like `provisioning_steps` inside of Laforge.

> Note that while filenames do **not** need to corraspond to the ID schema requirements, folder names in the path do. And while you can certainly name your .laforge files different from the ID defined in their block definition (and the source scripts for that matter), it is highly recommended you make them the same for easy identification.

The other global types can be referenced in this way:

```
NOTE: **/* represents recursive traversal of directories

/config/**/*      <~ competition-wide configuration definitions
/commands/**/*    <~ command definitions
/hosts/**/*       <~ host definitions
/networks/**/*    <~ network definitions
/identities/**/*  <~ identity definitions
/files/**/*       <~ remote_file definitions
/apps/**/*        <~ app definitions
/dns-records/**/* <~ dns_record definitions
```

Environment definitions (and their child definitions - `build`, `team`, and `provisioned_host`) exist in a slightly different format. Lets say you have an environment ID'd as "dev2018". It would have a fully qualified URI as `/envs/dev2018`. with a definition file located in at `/envs/dev2018/env.laforge`. It's children will live in subfolders, and are machine generated.

As an example, lets say we were using the Terraform Google Cloud Platform builder (ID=`tfgcp`) in `dev2018`. It's entire folder structure would look like this:

```
/envs/dev2018/
/envs/dev2018/env.laforge

/envs/dev2018/tfgcp/
/envs/dev2018/tfgcp/build.laforge
/envs/dev2018/tfgcp/state.db

/envs/dev2018/tfgcp/logs/
/envs/dev2018/tfgcp/logs/laforge.log

/envs/dev2018/tfgcp/data/
/envs/dev2018/tfgcp/data/ssh.pem
/envs/dev2018/tfgcp/data/ssh.pem.pub

# cache folder contains file uploads shared across teams.
/envs/dev2018/tfgcp/data/cache/
/envs/dev2018/tfgcp/data/cache/a9552c9e505df5893254bcac4e4044b6.zip

# teams are isolated from the build in a subfolder
/envs/dev2018/tfgcp/teams/

# TeamContext base directory
/envs/dev2018/tfgcp/teams/0/
/envs/dev2018/tfgcp/teams/0/team.laforge
/envs/dev2018/tfgcp/teams/0/infra.tf

# team-wide logs live in this directory, and are prefixed with a unix epoch.
/envs/dev2018/tfgcp/teams/0/logs/
/envs/dev2018/tfgcp/teams/0/logs/1234-terraform.stdout.log
/envs/dev2018/tfgcp/teams/0/logs/1234-terraform.stderr.log

# Provisioned hosts will exist in their own directory tree
/envs/dev2018/tfgcp/teams/0/networks/

# Each host is isolated into the network it's deployed into
/envs/dev2018/tfgcp/teams/0/networks/vdi/
/envs/dev2018/tfgcp/teams/0/networks/vdi/provisioned_network.laforge

# And then inside it's own folder
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/provisioned_host.laforge

# Each rendered asset for a host will live in the assets/ directory
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/assets/
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/assets/install-mysql.sh

# And has it's own log directory for activity about the host.
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/logs/
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/logs/1234-install-mysql.stdout.log
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/logs/1234-install-mysql.stderr.log

# And a directory for it's steps
/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/steps/0-install-mysql.laforge
```

Inside this build tree, everything is referenced by fully qualified path:

- **Environment ID**: `/envs/dev2018`
- **Build ID**: `/envs/dev2018/tfgcp`
- **Team ID**: `/envs/dev2018/tfgcp/teams/0`
- **Provisioned Network ID**: `/envs/dev2018/tfgcp/teams/0/networks/vdi`
- **Provisioned Host ID**: `/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01`
- **Provisioning Step ID**: `/envs/dev2018/tfgcp/teams/0/networks/vdi/hosts/ns01/steps/0-install-mysql`
