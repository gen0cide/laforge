# Laforge V3 Features

## Findings Tagging In Scripts

Allows developers to mark findings from within laforge scripts themselves using a common syntax.

### Syntax Example

```
#laforge:finding --name="" --description="" --severity=[1-10] --difficulty=[1-3]
```
These would be compiled as part of the final environment's output state.

## Terraform Remote State (Sharing)
Laforge Environments should share their terraform states. Not sure how we want to do this, but will figure something out.

## Universal Host & Network YAMLs
Host and Network YAMLs should reside within top level directories. Network definitions should be defined within the env.yml. 

## Terraform Top Level Directory
Move terraform states to a top level directory. 

## Multi-Level path referencing for scripts/files.
Scripts should be able to be pathed within subfolders for clarity.

## Replace remote provisoners with Laforge Custom provisioners
Instead of remote provisioning, allow for a more interactive user experience in build steps.

Examples:

 * Debugger
 * Local sleep
 * Reboot machine
 * Run local script
 
Features:

 * Logs output for each provisioning step uniquely
 * Allows debug output to not be piped directly to the console
 
## Remove environment prefixes for hostnames
This was only needed because of Splunk, but now that we have Elastic IPs we can easily just have a laforge mapping that maps IP to hostname and team ID.

## CDN Authenticated Links
Allow for environment building to generate precompiled links for files needed for retrieval - or perhaps a mechanism to pull them down internally. Need to flush this out.

## Laforge UI Server
A new subcommand `server` should be implemented to run an interactive WebUI to visualize Laforge functionality.

## Laforge Personal Configuration
`.lf_settings` should exist to allow the definition of callbacks, etc.

## Laforge Action Callbacks
Perform an arbitrary set of actions when something takes place.

## DNS Server Delta calculations
If genesis already exists, calculate the ns-updates needed to add or remove any host deltas that may exist.

## Create an environment specific DNS Configuration
With R53, we had the ability to define custom DNS records - we should make sure this happens again.

## Cross Team Command Execution
`laforge [ssh/powershell] --teams=[0-11] $hostname "command"`

## Laforge File Upload
`laforge scp $hostname $local_file $remote_file` (should work for windows as well)

## More laforge terraform subcommands

 * `laforge tf taint $objname`
 
## Laforge should handle signal traps gracefully
CTRL+C should do what it's expected to do and gracefully terminate terraform. Also timeouts for provisioners should be configurable in the host YAML.

## More comprehensive Scripting Guide
Need people to take more advantage of the scripting guidelines.

## Need better way to manage AMIs.
Need better way to manage the amis.json file.

## SSH keypair de-duplication
Only use a single SSH Key Pair object in AWS for each competition.

## Laforge Firewall Rewrite
Instead of de-duping public ports per team, firewall rules should be specified per host

## Genesis Host => Root DNS
This should just be renamed to what it is. Less confusing. And it should be treated 1:1 as a real host.

## OOB Network in each VPC
The DNS server, as well as potential jump hosts should live not inside VDI but inside a special management OOB network in each VPC.

## De-duplicate Terraform modules on host
There is no need to `terraform init` on every team. This is just stupid.

##
 