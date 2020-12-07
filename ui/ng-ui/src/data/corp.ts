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
  exposedTCPPorts: [
    '123',
    '135',
    '464',
    '49152-65535',
    '389',
    '636',
    '3268',
    '3269',
    '53',
    '88',
    '445',
    '3389'
  ],
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

const coins_head_01_provisioned: ProvisionedHost = {
  id: 'prov-CH01',
  subnetIP: '10.0.1.2',
  status: complete_status,
  provisionedNetwork: null,
  provisionedSteps: [],
  host: coins_heads_01
};

const dc_01_provisioned: ProvisionedHost = {
  id: 'prov-DC01',
  subnetIP: '10.0.1.10',
  status: in_progress_status,
  provisionedNetwork: null,
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
  provisionedHosts: [coins_head_01_provisioned, dc_01_provisioned],
  status: complete_status,
  network: corp_network,
  build: null
};

export { corp_network_provisioned };
