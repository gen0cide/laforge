import { ID, tagMap, varsMap } from './common.model';
import { Competition, Environment } from './environment.model';

export interface DNSRecord {
  id: ID;
  name: string;
  values: string[];
  type: string;
  zone: string;
  disabled: boolean;
  hcl_id?: string;
  vars?: varsMap[];
  tags: tagMap[];
  DNSRecordToEnvironment?: Environment;
}

export interface DNS {
  id: ID;
  type: string;
  root_domain: string;
  dns_servers: string[];
  ntp_servers: string[];
  hcl_id?: string;
  DNSToEnvironment?: Environment;
  DNSToCompetition?: Competition;
}
