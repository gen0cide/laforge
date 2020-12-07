import { Command } from './command.model';
import { ID, Status, Tag, User, varsMap } from './common.model';
import { DNSRecord } from './dns.model';
import { FileDelete, FileDownload, FileExtract } from './file.model';
import { ProvisionedNetwork } from './network.model';
import { Script } from './script.model';

interface Host {
  id: ID;
  hostname: string;
  OS: string;
  lastOctet: number;
  allowMacChanges: boolean;
  exposedTCPPorts: string[];
  exposedUDPPorts: string[];
  overridePassword: string;
  vars: varsMap[];
  userGroups: string[];
  dependsOn: Host[];
  maintainer: User;
  tags: Tag[];
  dnsRecords?: DNSRecord[];
  commands?: Command[];
  scripts?: Script[];
  fileDeletes?: FileDelete[];
  fileDownloads?: FileDownload[];
  fileExtracts?: FileExtract[];
}

interface ProvisionedHost {
  id: ID;
  subnetIP: string;
  status: Status;
  provisionedNetwork: ProvisionedNetwork;
  provisionedSteps: ProvisionedStep[];
  host: Host;
}

interface ProvisionedStep {
  id: ID;
  provisionType: string;
  stepNumber: number;
  provisionedHost: ProvisionedHost;
  status: Status;
  script?: Script;
  command?: Command;
  DNSRecord?: DNSRecord;
  fileDownload?: FileDownload;
  fileDelete?: FileDelete;
  fileExtract?: FileExtract;
}

export { Host, ProvisionedHost, ProvisionedStep };
