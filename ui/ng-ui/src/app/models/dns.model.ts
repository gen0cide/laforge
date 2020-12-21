import { configMap, ID, Tag, varsMap } from './common.model';

export interface DNSRecord {
  id: ID;
  name: string;
  values: string[];
  type: string;
  zone: string;
  vars: varsMap[];
  tags: Tag[];
  disabled: boolean;
}

export interface DNS {
  id: ID;
  type: string;
  rootDomain: string;
  DNSServers: string[];
  NTPServers: string[];
  config: configMap[];
}
