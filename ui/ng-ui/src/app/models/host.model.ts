import { AgentStatus } from './agent.model';
import { ID, Plan, Status, tagMap, varsMap } from './common.model';
import { Environment } from './environment.model';
import { ProvisionedNetwork } from './network.model';
import { ProvisioningStep } from './step.model';

export interface Disk {
  size: number;
  DiskToHost?: Host;
}

export interface Host {
  id: ID;
  hostname: string;
  OS: string;
  allow_mac_changes: boolean;
  exposed_tcp_ports: string[];
  exposed_udp_ports: string[];
  override_password: string;
  vars: varsMap[];
  user_groups: string[];
  tags: tagMap[];
  hcl_id?: string;
  description?: string;
  last_octet?: number;
  instance_size?: string;
  provision_steps?: string[];
  HostToDisk?: Disk;
  HostToEnvironment?: Environment;
}

export interface ProvisionedHost {
  id: ID;
  subnet_ip: string;
  ProvisionedHostToStatus: Status;
  combined_output?: string;
  ProvisionedHostToProvisionedNetwork?: ProvisionedNetwork;
  ProvisionedHostToHost?: Host;
  ProvisionedHostToProvisioningStep?: ProvisioningStep[];
  ProvisionedHostToAgentStatus?: AgentStatus;
  ProvisionedHostToPlan?: Plan;
}
