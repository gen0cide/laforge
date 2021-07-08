import { ProvisionStatus, Status, User } from 'src/app/models/common.model';
import { Network, ProvisionedNetwork } from 'src/app/models/network.model';
import { Host, ProvisionedHost } from '../app/models/host.model';

const default_user: User = {
  id: 'default',
  name: 'Bradley',
  uuid: 'QnJhZGxleQ==',
  email: 'harkerbd@gmail.com'
};

const complete_status: Status = {
  state: ProvisionStatus.ProvStatusComplete,
  startedAt: '1607366996285',
  endedAt: '1607377862255',
  failed: false,
  completed: true,
  error: ''
};

const in_progress_status: Status = {
  state: ProvisionStatus.ProvStatusInProgress,
  startedAt: '1607366996285',
  endedAt: '1607377862255',
  failed: false,
  completed: false,
  error: ''
};

const failed_status: Status = {
  state: ProvisionStatus.ProvStatusFailed,
  startedAt: '1607366996285',
  endedAt: '1607377862255',
  failed: true,
  completed: false,
  error: 'An unknown error occurred'
};

const coins_heads_01: Host = {
  id: 'CH01',
  hostname: 'heads-01',
  OS: 'ubuntu18',
  lastOctet: 113,
  allowMacChanges: false,
  exposedTCPPorts: ['22', '80', '443', '8546', '30301'],
  exposedUDPPorts: ['30303'],
  overridePassword: 'Preh!storic',
  vars: [
    {
      key: 'user_data_script_id',
      value: '/scripts/linux/ubuntu-userdata-script'
    }
  ],
  userGroups: [],
  dependsOn: [],
  maintainer: default_user,
  tags: [
    {
      id: 'focus-1',
      name: 'focus',
      description: 'services'
    }
  ]
};

const dc_01: Host = {
  id: 'DC01',
  hostname: 'corp-dc-01',
  OS: 'w2k16',
  lastOctet: 10,
  allowMacChanges: false,
  exposedTCPPorts: ['123', '135', '464', '49152-65535', '389', '636', '3268', '3269', '53', '88', '445', '3389'],
  exposedUDPPorts: ['53', '88', '389', '123', '464'],
  overridePassword: '',
  vars: [
    {
      key: 'force_master_dns',
      value: 'false'
    },
    {
      key: 'user_data_script_id',
      value: '/scripts/windows/powershell-userdata-script'
    }
  ],
  userGroups: [],
  dependsOn: [],
  maintainer: default_user,
  tags: [
    {
      id: 'focus-2',
      name: 'focus',
      description: 'full-access'
    },
    {
      id: 'network-1',
      name: 'network',
      description: 'corp'
    },
    {
      id: 'service-1',
      name: 'service',
      description: 'activedirectory'
    }
  ]
};

const coins_heads_01_provisioned: ProvisionedHost = {
  id: 'prov-CH01',
  subnetIP: '10.0.1.1',
  status: failed_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: coins_heads_01
};

const coins_heads_02_provisioned: ProvisionedHost = {
  id: 'prov-CH02',
  subnetIP: '10.0.1.2',
  status: failed_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: coins_heads_01
};

const coins_heads_03_provisioned: ProvisionedHost = {
  id: 'prov-CH03',
  subnetIP: '10.0.1.3',
  status: failed_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: coins_heads_01
};

const coins_heads_04_provisioned: ProvisionedHost = {
  id: 'prov-CH04',
  subnetIP: '10.0.1.4',
  status: failed_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: coins_heads_01
};

const coins_heads_05_provisioned: ProvisionedHost = {
  id: 'prov-CH05',
  subnetIP: '10.0.1.5',
  status: failed_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: coins_heads_01
};

const dc_01_provisioned: ProvisionedHost = {
  id: 'prov-DC01',
  subnetIP: '10.0.1.10',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_02_provisioned: ProvisionedHost = {
  id: 'prov-DC02',
  subnetIP: '10.0.1.12',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_03_provisioned: ProvisionedHost = {
  id: 'prov-DC03',
  subnetIP: '10.0.1.13',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_04_provisioned: ProvisionedHost = {
  id: 'prov-DC04',
  subnetIP: '10.0.1.14',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_05_provisioned: ProvisionedHost = {
  id: 'prov-DC05',
  subnetIP: '10.0.1.15',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_06_provisioned: ProvisionedHost = {
  id: 'prov-DC06',
  subnetIP: '10.0.1.16',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_07_provisioned: ProvisionedHost = {
  id: 'prov-DC07',
  subnetIP: '10.0.1.17',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_08_provisioned: ProvisionedHost = {
  id: 'prov-DC08',
  subnetIP: '10.0.1.18',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_09_provisioned: ProvisionedHost = {
  id: 'prov-DC09',
  subnetIP: '10.0.1.19',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const dc_10_provisioned: ProvisionedHost = {
  id: 'prov-DC10',
  subnetIP: '10.0.1.20',
  status: in_progress_status,
  // provisionedNetwork: null, // avoid circular dependencies
  provisionedSteps: [],
  host: dc_01
};

const corp_network: Network = {
  id: 'net-corp',
  name: 'corp',
  cidr: '10.0.1.0/24',
  vdiVisible: true,
  vars: [
    {
      key: 'authoritative_dns',
      value: 'true'
    },
    {
      key: 'authoritative_dns_ip',
      value: '10.0.1.10'
    },
    {
      key: 'authoritative_dns_host',
      value: 'corp-dc-01'
    },
    {
      key: 'ad_domain_dns',
      value: 'corp.dinobank.us'
    },
    {
      key: 'ad_domain_netbios',
      value: 'DINO'
    }
  ],
  tags: []
};

const corp_network_provisioned: ProvisionedNetwork = {
  id: 'prov-net-corp',
  name: 'corp01',
  cidr: '10.0.1.0/24',
  vars: [
    {
      key: 'authoritative_dns',
      value: 'true'
    },
    {
      key: 'authoritative_dns_ip',
      value: '10.0.1.10'
    },
    {
      key: 'authoritative_dns_host',
      value: 'corp-dc-01'
    },
    {
      key: 'ad_domain_dns',
      value: 'corp.dinobank.us'
    },
    {
      key: 'ad_domain_netbios',
      value: 'DINO'
    }
  ],
  tags: [],
  provisionedHosts: [
    coins_heads_01_provisioned,
    coins_heads_02_provisioned,
    coins_heads_03_provisioned,
    coins_heads_04_provisioned,
    coins_heads_05_provisioned,
    dc_01_provisioned,
    dc_02_provisioned,
    dc_03_provisioned,
    dc_04_provisioned,
    dc_05_provisioned,
    dc_06_provisioned,
    dc_07_provisioned,
    dc_08_provisioned,
    dc_09_provisioned,
    dc_10_provisioned
  ],
  status: complete_status,
  network: corp_network
  // build: null
};

export {
  coins_heads_01,
  dc_01,
  coins_heads_01_provisioned,
  dc_01_provisioned,
  corp_network,
  corp_network_provisioned,
  complete_status,
  in_progress_status,
  failed_status
};
