// import { Team, User } from 'src/app/models/common.model';
// import { DNS } from 'src/app/models/dns.model';
// import { Build, Competition, Environment } from 'src/app/models/environment.model';
// import { corp_network, coins_heads_01, dc_01, corp_network_provisioned } from './corp';

// const default_user: User = {
//   id: 'default',
//   name: 'Bradley',
//   uuid: 'QnJhZGxleQ==',
//   email: 'harkerbd@gmail.com'
// };
// const other_user: User = {
//   id: 'default',
//   name: 'Chike',
//   uuid: 'QnJhZGxleQ==',
//   email: 'chikeudenze@gmail.com'
// };

// const default_dns: DNS = {
//   id: 'default',
//   type: 'bind',
//   rootDomain: 'dinobank.us',
//   DNSServers: ['8.8.8.8', '8.8.4.4'],
//   NTPServers: ['129.6.15.28', '129.6.15.29'],
//   config: []
// };

// const cptc2019: Competition = {
//   id: 'cptc2019',
//   rootPassword: 'D!n0saur',
//   config: [],
//   dns: default_dns
// };

// const team1: Team = {
//   id: 'team1',
//   teamNumber: 1,
//   config: [],
//   revision: 1,
//   maintainer: default_user,
//   build: null, // Avoid circular dependencies
//   environment: null, // same as above
//   tags: [],
//   provisionedNetworks: [corp_network_provisioned]
// };

// const team2: Team = {
//   id: 'team2',
//   teamNumber: 2,
//   config: [],
//   revision: 1,
//   maintainer: default_user,
//   build: null, // Avoid circular dependencies
//   environment: null, // same as above
//   tags: [],
//   provisionedNetworks: [corp_network_provisioned]
// };

// const team3: Team = {
//   id: 'team3',
//   teamNumber: 3,
//   config: [],
//   revision: 1,
//   maintainer: default_user,
//   build: null, // Avoid circular dependencies
//   environment: null, // same as above
//   tags: [],
//   provisionedNetworks: [corp_network_provisioned]
// };

// const team4: Team = {
//   id: 'team4',
//   teamNumber: 4,
//   config: [],
//   revision: 1,
//   maintainer: default_user,
//   build: null, // Avoid circular dependencies
//   environment: null, // same as above
//   tags: [],
//   provisionedNetworks: [corp_network_provisioned]
// };

// const team5: Team = {
//   id: 'team5',
//   teamNumber: 5,
//   config: [],
//   revision: 1,
//   maintainer: default_user,
//   build: null, // Avoid circular dependencies
//   environment: null, // same as above
//   tags: [],
//   provisionedNetworks: [corp_network_provisioned]
// };

// const team6: Team = {
//   id: 'team6',
//   teamNumber: 6,
//   config: [],
//   revision: 1,
//   maintainer: default_user,
//   build: null, // Avoid circular dependencies
//   environment: null, // same as above
//   tags: [],
//   provisionedNetworks: [corp_network_provisioned]
// };

// const default_build: Build = {
//   id: 'bld-default',
//   revision: 1,
//   tags: [],
//   config: [],
//   maintainer: default_user,
//   teams: [team1, team2, team3, team4, team5, team6]
// };

// const bradley: Environment = {
//   id: 'bradley-1',
//   CompetitionID: 'cptc2019',
//   Name: 'bradley',
//   Description: 'Bradley Dev Env',
//   Builder: 'tfgcp',
//   TeamCount: 1,
//   AdminCIDRs: ['35.224.174.165/32', '35.193.160.100/32'],
//   ExposedVDIPorts: [],
//   tags: [],
//   config: [
//     {
//       key: 'vpc_cidr',
//       value: '10.0.0.0/8'
//     },
//     {
//       key: 'admin_ip',
//       value: '35.224.174.165'
//     },
//     {
//       key: 'vdi_whitelist',
//       value: '0.0.0.0/0'
//     },
//     {
//       key: 'gcp_cred_file',
//       value: '/cptc/auth/infra-test-environment.json'
//     },
//     {
//       key: 'gcp_project',
//       value: 'infra-test-environment'
//     },
//     {
//       key: 'gcp_region',
//       value: 'us-central1'
//     },
//     {
//       key: 'gcp_storage_bucket',
//       value: 'us-central1-a'
//     },
//     {
//       key: 'gcp_dns_zone_id',
//       value: 'dinobank-us-public'
//     },
//     {
//       key: 'gcp_project_ssh_pubkey',
//       value: 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKjDGacIV4OKVlghKlhDIVueqyrJHalanMF9gAeh+OOp root'
//     },
//     {
//       key: 'etcd_username',
//       value: 'root'
//     },
//     {
//       key: 'etcd_password',
//       value: 'FAXGLAPGOYWEZMPZ'
//     },
//     {
//       key: 'etcd_master',
//       value: 'portal668-0.cptc2019.577934193.composedb.com:19452'
//     },
//     {
//       key: 'etcd_slave',
//       value: 'portal496-1.cptc2019.577934193.composedb.com:19452'
//     },
//     {
//       key: 'master_dns_server',
//       value: '10.0.254.10'
//     }
//   ],
//   maintainer: default_user,
//   networks: [corp_network],
//   hosts: [coins_heads_01, dc_01],
//   build: default_build,
//   competition: cptc2019
// };
// const chike: Environment = {
//   id: 'chike-1',
//   CompetitionID: 'cptc2019',
//   Name: 'Chike',
//   Description: 'Chike Dev Env',
//   Builder: 'tfgcp',
//   TeamCount: 1,
//   AdminCIDRs: ['35.224.174.165/32', '35.193.160.100/32'],
//   ExposedVDIPorts: [],
//   tags: [],
//   config: [
//     {
//       key: 'vpc_cidr',
//       value: '10.0.0.0/8'
//     },
//     {
//       key: 'admin_ip',
//       value: '35.224.174.165'
//     },
//     {
//       key: 'vdi_whitelist',
//       value: '0.0.0.0/0'
//     },
//     {
//       key: 'gcp_cred_file',
//       value: '/cptc/auth/infra-test-environment.json'
//     },
//     {
//       key: 'gcp_project',
//       value: 'infra-test-environment'
//     },
//     {
//       key: 'gcp_region',
//       value: 'us-central1'
//     },
//     {
//       key: 'gcp_storage_bucket',
//       value: 'us-central1-a'
//     },
//     {
//       key: 'gcp_dns_zone_id',
//       value: 'dinobank-us-public'
//     },
//     {
//       key: 'gcp_project_ssh_pubkey',
//       value: 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKjDGacIV4OKVlghKlhDIVueqyrJHalanMF9gAeh+OOp root'
//     },
//     {
//       key: 'etcd_username',
//       value: 'root'
//     },
//     {
//       key: 'etcd_password',
//       value: 'FAXGLAPGOYWEZMPZ'
//     },
//     {
//       key: 'etcd_master',
//       value: 'portal668-0.cptc2019.577934193.composedb.com:19452'
//     },
//     {
//       key: 'etcd_slave',
//       value: 'portal496-1.cptc2019.577934193.composedb.com:19452'
//     },
//     {
//       key: 'master_dns_server',
//       value: '10.0.254.10'
//     }
//   ],
//   maintainer: other_user,
//   networks: [corp_network],
//   hosts: [coins_heads_01, dc_01],
//   build: default_build,
//   competition: cptc2019
// };

// export { bradley, chike };
