import { configMap, ID, Tag, varsMap } from './common.model';

export interface DNSRecord {
  id: ID;
  hcl_id?: string;
  name: string;
  values: string[];
  type: string;
  zone: string;
  vars: varsMap[];
  dnsRecordToTag: Tag[];
  disabled: boolean;
}

export interface DNS {
  id: ID;
  type: string;
  rootDomain: string;
  DNSServers: string[];
  NTPServer: string[];
  config: configMap[];
}
