import { configMap, ID, Tag, varsMap } from './common.model';

interface DNSRecord {
  id: ID;
  name: String;
  values: String[];
  type: String;
  zone: String;
  vars: varsMap[];
  tags: Tag[];
  disabled: Boolean;
}

interface DNS {
  id: ID;
  type: String;
  rootDomain: String;
  DNSServers: String[];
  NTPServer: String[];
  config: configMap[];
}

export {
  DNSRecord,
  DNS
}