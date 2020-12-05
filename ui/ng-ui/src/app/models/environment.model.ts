import { TemplateAst } from '@angular/compiler'
import { ID, Tag, configMap, User, Team } from './common.model'
import { DNS, DNSRecord } from './dns.model'
import { Host } from './host.model'
import { Network } from './network.model'

interface Build {
  id: ID;
  revision: Number;
  tags: Tag[];
  config: configMap[];
  maintainer: User;
  teams: Team[];
}

interface Competition {
  id: ID;
  rootPassword: String;
  config: configMap[];
  dns: DNS;
}

interface Environment {
  id: ID;
  CompetitionID: String;
  Name: String;
  Description: String;
  Builder: String;
  TeamCount: Number;
  AdminCIDRs: String[];
  ExposedVDIPorts: String[];
  tags: Tag[];
  config: configMap[];
  maintainer: User;
  networks: Network[];
  hosts: Host[]
  build: Build;
  competition: Competition;
}

export {
  Environment,
  Build
}