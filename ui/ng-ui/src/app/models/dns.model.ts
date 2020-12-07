import { configMap, ID, Tag, varsMap } from './common.model';

interface DNSRecord {
  id: ID;
  name: string;
  values: string[];
  type: string;
  zone: string;
  vars: varsMap[];
  tags: Tag[];
  disabled: boolean;
}

interface DNS {
  id: ID;
  type: string;
  rootDomain: string;
  DNSServers: string[];
  NTPServers: string[];
  config: configMap[];
}

export { DNSRecord, DNS };
