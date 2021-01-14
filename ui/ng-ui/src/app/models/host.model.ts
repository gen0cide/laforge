import { AgentStatus } from './agent.model';
import { Command } from './command.model';
import { ID, Status, Tag, User, varsMap } from './common.model';
import { DNSRecord } from './dns.model';
import { FileDelete, FileDownload, FileExtract, RemoteFile } from './file.model';
// import { ProvisionedNetwork } from './network.model';
import { Script } from './script.model';

export interface Host {
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

export interface ProvisionedHost {
  id: ID;
  subnetIP: string;
  status?: Status;
  // provisionedNetwork: ProvisionedNetwork; * avoids circular dependencies *
  provisionedSteps?: ProvisionedStep[] /* optional, allows us to query steps later */;
  host: Host;
  heartbeat?: AgentStatus;
}

export interface ProvisionedStep {
  id: ID;
  provisionType: string;
  stepNumber: number;
  provisionedHost: ProvisionedHost;
  status?: Status;
  script?: Script;
  command?: Command;
  DNSRecord?: DNSRecord;
  remoteFile?: RemoteFile;
}
