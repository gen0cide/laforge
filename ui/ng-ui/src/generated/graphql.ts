import { Injectable } from '@angular/core';
import { gql } from 'apollo-angular';
import * as Apollo from 'apollo-angular';
export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Time: any;
};

export enum LaForgeAgentCommand {
  Default = 'DEFAULT',
  Delete = 'DELETE',
  Reboot = 'REBOOT',
  Extract = 'EXTRACT',
  Download = 'DOWNLOAD',
  Createuser = 'CREATEUSER',
  Createuserpass = 'CREATEUSERPASS',
  Addtogroup = 'ADDTOGROUP',
  Execute = 'EXECUTE',
  Validate = 'VALIDATE',
  Changeperms = 'CHANGEPERMS',
  Appendfile = 'APPENDFILE'
}

export type LaForgeAgentStatus = {
  __typename?: 'AgentStatus';
  clientId: Scalars['String'];
  hostname: Scalars['String'];
  upTime: Scalars['Int'];
  bootTime: Scalars['Int'];
  numProcs: Scalars['Int'];
  OS: Scalars['String'];
  hostID: Scalars['String'];
  load1?: Maybe<Scalars['Float']>;
  load5?: Maybe<Scalars['Float']>;
  load15?: Maybe<Scalars['Float']>;
  totalMem: Scalars['Int'];
  freeMem: Scalars['Int'];
  usedMem: Scalars['Int'];
  timestamp: Scalars['Int'];
};

export type LaForgeAgentStatusBatch = {
  __typename?: 'AgentStatusBatch';
  agentStatuses: Array<Maybe<LaForgeAgentStatus>>;
  pageInfo: LaForgeLaForgePageInfo;
};

export type LaForgeAgentTask = {
  __typename?: 'AgentTask';
  id: Scalars['ID'];
  args?: Maybe<Scalars['String']>;
  command: LaForgeAgentCommand;
  number: Scalars['Int'];
  output?: Maybe<Scalars['String']>;
  state: LaForgeAgentTaskState;
  error_message?: Maybe<Scalars['String']>;
};

export enum LaForgeAgentTaskState {
  Awaiting = 'AWAITING',
  Inprogress = 'INPROGRESS',
  Failed = 'FAILED',
  Complete = 'COMPLETE'
}

export type LaForgeAuthUser = {
  __typename?: 'AuthUser';
  id: Scalars['ID'];
  username: Scalars['String'];
  role: LaForgeRoleLevel;
  provider: LaForgeProviderType;
  first_name: Scalars['String'];
  last_name: Scalars['String'];
  email: Scalars['String'];
  phone: Scalars['String'];
  company: Scalars['String'];
  occupation: Scalars['String'];
  publicKey: Scalars['String'];
};

export type LaForgeBuild = {
  __typename?: 'Build';
  id: Scalars['ID'];
  revision: Scalars['Int'];
  environment_revision: Scalars['Int'];
  completed_plan: Scalars['Boolean'];
  buildToStatus: LaForgeStatus;
  buildToEnvironment: LaForgeEnvironment;
  buildToCompetition: LaForgeCompetition;
  buildToProvisionedNetwork: Array<Maybe<LaForgeProvisionedNetwork>>;
  buildToTeam: Array<Maybe<LaForgeTeam>>;
  buildToPlan: Array<Maybe<LaForgePlan>>;
  BuildToLatestBuildCommit?: Maybe<LaForgeBuildCommit>;
  BuildToBuildCommits: Array<Maybe<LaForgeBuildCommit>>;
};

export type LaForgeBuildCommit = {
  __typename?: 'BuildCommit';
  id: Scalars['ID'];
  type: LaForgeBuildCommitType;
  revision: Scalars['Int'];
  state: LaForgeBuildCommitState;
  BuildCommitToBuild: LaForgeBuild;
  BuildCommitToPlanDiffs: Array<Maybe<LaForgePlanDiff>>;
};

export enum LaForgeBuildCommitState {
  Planning = 'PLANNING',
  Inprogress = 'INPROGRESS',
  Applied = 'APPLIED',
  Cancelled = 'CANCELLED',
  Approved = 'APPROVED'
}

export enum LaForgeBuildCommitType {
  Root = 'ROOT',
  Rebuild = 'REBUILD',
  Delete = 'DELETE'
}

export type LaForgeCommand = {
  __typename?: 'Command';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  name: Scalars['String'];
  description: Scalars['String'];
  program: Scalars['String'];
  args: Array<Maybe<Scalars['String']>>;
  ignoreErrors: Scalars['Boolean'];
  disabled: Scalars['Boolean'];
  cooldown: Scalars['Int'];
  timeout: Scalars['Int'];
  vars?: Maybe<Array<Maybe<LaForgeVarsMap>>>;
  tags?: Maybe<Array<Maybe<LaForgeTagMap>>>;
  CommandToEnvironment: LaForgeEnvironment;
};

export type LaForgeCompetition = {
  __typename?: 'Competition';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  root_password: Scalars['String'];
  config?: Maybe<Array<Maybe<LaForgeConfigMap>>>;
  tags?: Maybe<Array<Maybe<LaForgeTagMap>>>;
  competitionToDNS: Array<Maybe<LaForgeDns>>;
  CompetitionToEnvironment: LaForgeEnvironment;
  CompetitionToBuild: Array<Maybe<LaForgeBuild>>;
};

export type LaForgeDns = {
  __typename?: 'DNS';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  type: Scalars['String'];
  root_domain: Scalars['String'];
  dns_servers: Array<Maybe<Scalars['String']>>;
  ntp_servers: Array<Maybe<Scalars['String']>>;
  config?: Maybe<Array<Maybe<LaForgeConfigMap>>>;
  DNSToEnvironment: Array<Maybe<LaForgeEnvironment>>;
  DNSToCompetition: Array<Maybe<LaForgeCompetition>>;
};

export type LaForgeDnsRecord = {
  __typename?: 'DNSRecord';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  name: Scalars['String'];
  values: Array<Maybe<Scalars['String']>>;
  type: Scalars['String'];
  zone: Scalars['String'];
  vars: Array<Maybe<LaForgeVarsMap>>;
  disabled: Scalars['Boolean'];
  tags: Array<Maybe<LaForgeTagMap>>;
  DNSRecordToEnvironment: LaForgeEnvironment;
};

export type LaForgeDisk = {
  __typename?: 'Disk';
  size: Scalars['Int'];
  DiskToHost: LaForgeHost;
};

export type LaForgeEnvironment = {
  __typename?: 'Environment';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  competition_id: Scalars['String'];
  name: Scalars['String'];
  description: Scalars['String'];
  builder: Scalars['String'];
  team_count: Scalars['Int'];
  revision: Scalars['Int'];
  admin_cidrs: Array<Maybe<Scalars['String']>>;
  exposed_vdi_ports: Array<Maybe<Scalars['String']>>;
  config?: Maybe<Array<Maybe<LaForgeConfigMap>>>;
  tags?: Maybe<Array<Maybe<LaForgeTagMap>>>;
  EnvironmentToUser: Array<Maybe<LaForgeUser>>;
  EnvironmentToHost: Array<Maybe<LaForgeHost>>;
  EnvironmentToCompetition: Array<Maybe<LaForgeCompetition>>;
  EnvironmentToIdentity: Array<Maybe<LaForgeIdentity>>;
  EnvironmentToCommand: Array<Maybe<LaForgeCommand>>;
  EnvironmentToScript: Array<Maybe<LaForgeScript>>;
  EnvironmentToFileDownload: Array<Maybe<LaForgeFileDownload>>;
  EnvironmentToFileDelete: Array<Maybe<LaForgeFileDelete>>;
  EnvironmentToFileExtract: Array<Maybe<LaForgeFileExtract>>;
  EnvironmentToDNSRecord: Array<Maybe<LaForgeDnsRecord>>;
  EnvironmentToDNS: Array<Maybe<LaForgeDns>>;
  EnvironmentToNetwork: Array<Maybe<LaForgeNetwork>>;
  EnvironmentToBuild: Array<Maybe<LaForgeBuild>>;
  EnvironmentToRepository: Array<Maybe<LaForgeRepository>>;
};

export type LaForgeFileDelete = {
  __typename?: 'FileDelete';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  path: Scalars['String'];
  tags: Array<Maybe<LaForgeTagMap>>;
  FileDeleteToEnvironment: LaForgeEnvironment;
};

export type LaForgeFileDownload = {
  __typename?: 'FileDownload';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  sourceType: Scalars['String'];
  source: Scalars['String'];
  destination: Scalars['String'];
  template: Scalars['Boolean'];
  perms: Scalars['String'];
  disabled: Scalars['Boolean'];
  md5: Scalars['String'];
  absPath: Scalars['String'];
  tags: Array<Maybe<LaForgeTagMap>>;
  FileDownloadToEnvironment: LaForgeEnvironment;
};

export type LaForgeFileExtract = {
  __typename?: 'FileExtract';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  source: Scalars['String'];
  destination: Scalars['String'];
  type: Scalars['String'];
  tags: Array<Maybe<LaForgeTagMap>>;
  FileExtractToEnvironment: LaForgeEnvironment;
};

export type LaForgeFinding = {
  __typename?: 'Finding';
  name: Scalars['String'];
  description: Scalars['String'];
  severity: LaForgeFindingSeverity;
  difficulty: LaForgeFindingDifficulty;
  tags: Array<Maybe<LaForgeTagMap>>;
  FindingToUser: Array<Maybe<LaForgeUser>>;
  FindingToScript: LaForgeScript;
  FindingToEnvironment: LaForgeEnvironment;
};

export enum LaForgeFindingDifficulty {
  ZeroDifficulty = 'ZeroDifficulty',
  NoviceDifficulty = 'NoviceDifficulty',
  AdvancedDifficulty = 'AdvancedDifficulty',
  ExpertDifficulty = 'ExpertDifficulty',
  NullDifficulty = 'NullDifficulty'
}

export enum LaForgeFindingSeverity {
  ZeroSeverity = 'ZeroSeverity',
  LowSeverity = 'LowSeverity',
  MediumSeverity = 'MediumSeverity',
  HighSeverity = 'HighSeverity',
  CriticalSeverity = 'CriticalSeverity',
  NullSeverity = 'NullSeverity'
}

export type LaForgeHost = {
  __typename?: 'Host';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  hostname: Scalars['String'];
  description: Scalars['String'];
  OS: Scalars['String'];
  last_octet: Scalars['Int'];
  instance_size: Scalars['String'];
  allow_mac_changes: Scalars['Boolean'];
  exposed_tcp_ports: Array<Maybe<Scalars['String']>>;
  exposed_udp_ports: Array<Maybe<Scalars['String']>>;
  override_password: Scalars['String'];
  vars?: Maybe<Array<Maybe<LaForgeVarsMap>>>;
  user_groups: Array<Maybe<Scalars['String']>>;
  provision_steps: Array<Maybe<Scalars['String']>>;
  tags: Array<Maybe<LaForgeTagMap>>;
  HostToDisk: LaForgeDisk;
  HostToEnvironment: LaForgeEnvironment;
};

export type LaForgeIdentity = {
  __typename?: 'Identity';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  first_name: Scalars['String'];
  last_name: Scalars['String'];
  email: Scalars['String'];
  password: Scalars['String'];
  description: Scalars['String'];
  avatar_file: Scalars['String'];
  vars: Array<Maybe<LaForgeVarsMap>>;
  tags: Array<Maybe<LaForgeTagMap>>;
  IdentityToEnvironment: LaForgeEnvironment;
};

export type LaForgeLaForgePageInfo = {
  __typename?: 'LaForgePageInfo';
  total: Scalars['Int'];
  nextOffset: Scalars['Int'];
};

export type LaForgeMutation = {
  __typename?: 'Mutation';
  loadEnvironment?: Maybe<Array<Maybe<LaForgeEnvironment>>>;
  createBuild?: Maybe<LaForgeBuild>;
  deleteUser: Scalars['Boolean'];
  executePlan?: Maybe<LaForgeBuild>;
  deleteBuild: Scalars['Boolean'];
  createTask: Scalars['Boolean'];
  rebuild: Scalars['Boolean'];
  approveCommit: Scalars['Boolean'];
  cancelCommit: Scalars['Boolean'];
  createEnviromentFromRepo: Array<Maybe<LaForgeEnvironment>>;
  updateEnviromentViaPull: Array<Maybe<LaForgeEnvironment>>;
  modifySelfPassword: Scalars['Boolean'];
  modifySelfUserInfo?: Maybe<LaForgeAuthUser>;
  createUser?: Maybe<LaForgeAuthUser>;
  modifyAdminUserInfo?: Maybe<LaForgeAuthUser>;
  modifyAdminPassword: Scalars['Boolean'];
};

export type LaForgeMutationLoadEnvironmentArgs = {
  envFilePath: Scalars['String'];
};

export type LaForgeMutationCreateBuildArgs = {
  envUUID: Scalars['String'];
  renderFiles?: Scalars['Boolean'];
};

export type LaForgeMutationDeleteUserArgs = {
  userUUID: Scalars['String'];
};

export type LaForgeMutationExecutePlanArgs = {
  buildUUID: Scalars['String'];
};

export type LaForgeMutationDeleteBuildArgs = {
  buildUUID: Scalars['String'];
};

export type LaForgeMutationCreateTaskArgs = {
  proHostUUID: Scalars['String'];
  command: LaForgeAgentCommand;
  args: Scalars['String'];
};

export type LaForgeMutationRebuildArgs = {
  rootPlans: Array<Maybe<Scalars['String']>>;
};

export type LaForgeMutationApproveCommitArgs = {
  commitUUID: Scalars['String'];
};

export type LaForgeMutationCancelCommitArgs = {
  commitUUID: Scalars['String'];
};

export type LaForgeMutationCreateEnviromentFromRepoArgs = {
  repoURL: Scalars['String'];
  branchName?: Scalars['String'];
  repoName?: Scalars['String'];
  envFilePath: Scalars['String'];
};

export type LaForgeMutationUpdateEnviromentViaPullArgs = {
  envUUID: Scalars['String'];
};

export type LaForgeMutationModifySelfPasswordArgs = {
  currentPassword: Scalars['String'];
  newPassword: Scalars['String'];
};

export type LaForgeMutationModifySelfUserInfoArgs = {
  firstName?: Maybe<Scalars['String']>;
  lastName?: Maybe<Scalars['String']>;
  email?: Maybe<Scalars['String']>;
  phone?: Maybe<Scalars['String']>;
  company?: Maybe<Scalars['String']>;
  occupation?: Maybe<Scalars['String']>;
};

export type LaForgeMutationCreateUserArgs = {
  username: Scalars['String'];
  password: Scalars['String'];
  role: LaForgeRoleLevel;
  provider: LaForgeProviderType;
};

export type LaForgeMutationModifyAdminUserInfoArgs = {
  userID: Scalars['String'];
  username?: Maybe<Scalars['String']>;
  firstName?: Maybe<Scalars['String']>;
  lastName?: Maybe<Scalars['String']>;
  email?: Maybe<Scalars['String']>;
  phone?: Maybe<Scalars['String']>;
  company?: Maybe<Scalars['String']>;
  occupation?: Maybe<Scalars['String']>;
  role?: Maybe<LaForgeRoleLevel>;
  provider?: Maybe<LaForgeProviderType>;
};

export type LaForgeMutationModifyAdminPasswordArgs = {
  userID: Scalars['String'];
  newPassword: Scalars['String'];
};

export type LaForgeNetwork = {
  __typename?: 'Network';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  name: Scalars['String'];
  cidr: Scalars['String'];
  vdi_visible: Scalars['Boolean'];
  vars?: Maybe<Array<Maybe<LaForgeVarsMap>>>;
  tags: Array<Maybe<LaForgeTagMap>>;
  NetworkToEnvironment: LaForgeEnvironment;
};

export type LaForgePlan = {
  __typename?: 'Plan';
  id: Scalars['ID'];
  step_number: Scalars['Int'];
  type: LaForgePlanType;
  build_id: Scalars['String'];
  NextPlan: Array<Maybe<LaForgePlan>>;
  PrevPlan: Array<Maybe<LaForgePlan>>;
  PlanToBuild: LaForgeBuild;
  PlanToTeam: LaForgeTeam;
  PlanToProvisionedNetwork: LaForgeProvisionedNetwork;
  PlanToProvisionedHost: LaForgeProvisionedHost;
  PlanToProvisioningStep: LaForgeProvisioningStep;
  PlanToStatus: LaForgeStatus;
  PlanToPlanDiffs: Array<Maybe<LaForgePlanDiff>>;
};

export type LaForgePlanDiff = {
  __typename?: 'PlanDiff';
  id: Scalars['ID'];
  revision: Scalars['Int'];
  new_state: LaForgeProvisionStatus;
  PlanDiffToBuildCommit: LaForgeBuildCommit;
  PlanDiffToPlan: LaForgePlan;
};

export enum LaForgePlanType {
  StartBuild = 'start_build',
  StartTeam = 'start_team',
  ProvisionNetwork = 'provision_network',
  ProvisionHost = 'provision_host',
  ExecuteStep = 'execute_step',
  Undefined = 'undefined'
}

export enum LaForgeProviderType {
  Local = 'LOCAL',
  Github = 'GITHUB',
  Openid = 'OPENID',
  Undefined = 'UNDEFINED'
}

export enum LaForgeProvisionStatus {
  Planning = 'PLANNING',
  Awaiting = 'AWAITING',
  Inprogress = 'INPROGRESS',
  Failed = 'FAILED',
  Complete = 'COMPLETE',
  Tainted = 'TAINTED',
  Undefined = 'UNDEFINED',
  Todelete = 'TODELETE',
  Deleteinprogress = 'DELETEINPROGRESS',
  Deleted = 'DELETED',
  Torebuild = 'TOREBUILD'
}

export enum LaForgeProvisionStatusFor {
  Build = 'Build',
  Team = 'Team',
  Plan = 'Plan',
  ProvisionedNetwork = 'ProvisionedNetwork',
  ProvisionedHost = 'ProvisionedHost',
  ProvisioningStep = 'ProvisioningStep',
  Undefined = 'Undefined'
}

export type LaForgeProvisionedHost = {
  __typename?: 'ProvisionedHost';
  id: Scalars['ID'];
  subnet_ip: Scalars['String'];
  ProvisionedHostToStatus: LaForgeStatus;
  ProvisionedHostToProvisionedNetwork: LaForgeProvisionedNetwork;
  ProvisionedHostToHost: LaForgeHost;
  ProvisionedHostToProvisioningStep: Array<Maybe<LaForgeProvisioningStep>>;
  ProvisionedHostToAgentStatus?: Maybe<LaForgeAgentStatus>;
  ProvisionedHostToPlan: LaForgePlan;
};

export type LaForgeProvisionedNetwork = {
  __typename?: 'ProvisionedNetwork';
  id: Scalars['ID'];
  name: Scalars['String'];
  cidr: Scalars['String'];
  ProvisionedNetworkToStatus: LaForgeStatus;
  ProvisionedNetworkToNetwork: LaForgeNetwork;
  ProvisionedNetworkToBuild: LaForgeBuild;
  ProvisionedNetworkToTeam: LaForgeTeam;
  ProvisionedNetworkToProvisionedHost: Array<Maybe<LaForgeProvisionedHost>>;
  ProvisionedNetworkToPlan: LaForgePlan;
};

export type LaForgeProvisioningStep = {
  __typename?: 'ProvisioningStep';
  id: Scalars['ID'];
  type: LaForgeProvisioningStepType;
  step_number: Scalars['Int'];
  ProvisioningStepToStatus: LaForgeStatus;
  ProvisioningStepToProvisionedHost: LaForgeProvisionedHost;
  ProvisioningStepToScript?: Maybe<LaForgeScript>;
  ProvisioningStepToCommand?: Maybe<LaForgeCommand>;
  ProvisioningStepToDNSRecord?: Maybe<LaForgeDnsRecord>;
  ProvisioningStepToFileDelete?: Maybe<LaForgeFileDelete>;
  ProvisioningStepToFileDownload?: Maybe<LaForgeFileDownload>;
  ProvisioningStepToFileExtract?: Maybe<LaForgeFileExtract>;
  ProvisioningStepToPlan?: Maybe<LaForgePlan>;
};

export enum LaForgeProvisioningStepType {
  Script = 'Script',
  Command = 'Command',
  DnsRecord = 'DNSRecord',
  FileDelete = 'FileDelete',
  FileDownload = 'FileDownload',
  FileExtract = 'FileExtract',
  Undefined = 'Undefined'
}

export type LaForgeQuery = {
  __typename?: 'Query';
  environments?: Maybe<Array<Maybe<LaForgeEnvironment>>>;
  environment?: Maybe<LaForgeEnvironment>;
  provisionedHost?: Maybe<LaForgeProvisionedHost>;
  provisionedNetwork?: Maybe<LaForgeProvisionedNetwork>;
  provisionedStep?: Maybe<LaForgeProvisioningStep>;
  plan?: Maybe<LaForgePlan>;
  build?: Maybe<LaForgeBuild>;
  status?: Maybe<LaForgeStatus>;
  agentStatus?: Maybe<LaForgeAgentStatus>;
  getServerTasks?: Maybe<Array<Maybe<LaForgeServerTask>>>;
  currentUser?: Maybe<LaForgeAuthUser>;
  getUserList?: Maybe<Array<Maybe<LaForgeAuthUser>>>;
  getCurrentUserTasks?: Maybe<Array<Maybe<LaForgeServerTask>>>;
  getAgentTasks?: Maybe<Array<Maybe<LaForgeAgentTask>>>;
  getAllAgentStatus?: Maybe<LaForgeAgentStatusBatch>;
  getAllPlanStatus?: Maybe<LaForgeStatusBatch>;
};

export type LaForgeQueryEnvironmentArgs = {
  envUUID: Scalars['String'];
};

export type LaForgeQueryProvisionedHostArgs = {
  proHostUUID: Scalars['String'];
};

export type LaForgeQueryProvisionedNetworkArgs = {
  proNetUUID: Scalars['String'];
};

export type LaForgeQueryProvisionedStepArgs = {
  proStepUUID: Scalars['String'];
};

export type LaForgeQueryPlanArgs = {
  planUUID: Scalars['String'];
};

export type LaForgeQueryBuildArgs = {
  buildUUID: Scalars['String'];
};

export type LaForgeQueryStatusArgs = {
  statusUUID: Scalars['String'];
};

export type LaForgeQueryAgentStatusArgs = {
  clientId: Scalars['String'];
};

export type LaForgeQueryGetAgentTasksArgs = {
  proStepUUID: Scalars['String'];
};

export type LaForgeQueryGetAllAgentStatusArgs = {
  buildUUID: Scalars['String'];
  count: Scalars['Int'];
  offset: Scalars['Int'];
};

export type LaForgeQueryGetAllPlanStatusArgs = {
  buildUUID: Scalars['String'];
  count: Scalars['Int'];
  offset: Scalars['Int'];
};

export type LaForgeRepository = {
  __typename?: 'Repository';
  id: Scalars['ID'];
  repo_url: Scalars['String'];
  branch_name: Scalars['String'];
  environment_filepath: Scalars['String'];
  commit_info: Scalars['String'];
};

export enum LaForgeRoleLevel {
  Admin = 'ADMIN',
  User = 'USER',
  Undefined = 'UNDEFINED'
}

export type LaForgeScript = {
  __typename?: 'Script';
  id: Scalars['ID'];
  hcl_id: Scalars['String'];
  name: Scalars['String'];
  language: Scalars['String'];
  description: Scalars['String'];
  source: Scalars['String'];
  source_type: Scalars['String'];
  cooldown: Scalars['Int'];
  timeout: Scalars['Int'];
  ignore_errors: Scalars['Boolean'];
  args: Array<Maybe<Scalars['String']>>;
  disabled: Scalars['Boolean'];
  vars?: Maybe<Array<Maybe<LaForgeVarsMap>>>;
  absPath: Scalars['String'];
  tags?: Maybe<Array<Maybe<LaForgeTagMap>>>;
  scriptToFinding: Array<Maybe<LaForgeFinding>>;
  ScriptToEnvironment: LaForgeEnvironment;
};

export type LaForgeServerTask = {
  __typename?: 'ServerTask';
  id: Scalars['ID'];
  type: LaForgeServerTaskType;
  start_time?: Maybe<Scalars['Time']>;
  end_time?: Maybe<Scalars['Time']>;
  errors?: Maybe<Array<Maybe<Scalars['String']>>>;
  log_file_path?: Maybe<Scalars['String']>;
  ServerTaskToAuthUser: LaForgeAuthUser;
  ServerTaskToStatus: LaForgeStatus;
  ServerTaskToEnvironment?: Maybe<LaForgeEnvironment>;
  ServerTaskToBuild?: Maybe<LaForgeBuild>;
};

export enum LaForgeServerTaskType {
  Loadenv = 'LOADENV',
  Createbuild = 'CREATEBUILD',
  Renderfiles = 'RENDERFILES',
  Deletebuild = 'DELETEBUILD',
  Rebuild = 'REBUILD',
  Executebuild = 'EXECUTEBUILD'
}

export type LaForgeStatus = {
  __typename?: 'Status';
  id: Scalars['ID'];
  state: LaForgeProvisionStatus;
  status_for: LaForgeProvisionStatusFor;
  started_at: Scalars['String'];
  ended_at: Scalars['String'];
  failed: Scalars['Boolean'];
  completed: Scalars['Boolean'];
  error?: Maybe<Scalars['String']>;
};

export type LaForgeStatusBatch = {
  __typename?: 'StatusBatch';
  statuses: Array<Maybe<LaForgeStatus>>;
  pageInfo: LaForgeLaForgePageInfo;
};

export type LaForgeSubscription = {
  __typename?: 'Subscription';
  updatedAgentStatus: LaForgeAgentStatus;
  updatedStatus: LaForgeStatus;
  updatedServerTask: LaForgeServerTask;
  updatedBuild: LaForgeBuild;
  updatedCommit: LaForgeBuildCommit;
  updatedAgentTask: LaForgeAgentTask;
};

export type LaForgeTeam = {
  __typename?: 'Team';
  id: Scalars['ID'];
  team_number: Scalars['Int'];
  TeamToBuild: LaForgeBuild;
  TeamToStatus: LaForgeStatus;
  TeamToProvisionedNetwork: Array<Maybe<LaForgeProvisionedNetwork>>;
  TeamToPlan: LaForgePlan;
};

export type LaForgeUser = {
  __typename?: 'User';
  id: Scalars['ID'];
  name: Scalars['String'];
  uuid: Scalars['String'];
  email: Scalars['String'];
};

export type LaForgeConfigMap = {
  __typename?: 'configMap';
  key: Scalars['String'];
  value: Scalars['String'];
};

export type LaForgeTagMap = {
  __typename?: 'tagMap';
  key: Scalars['String'];
  value: Scalars['String'];
};

export type LaForgeVarsMap = {
  __typename?: 'varsMap';
  key: Scalars['String'];
  value: Scalars['String'];
};

export type LaForgeGetUserListQueryVariables = Exact<{ [key: string]: never }>;

export type LaForgeGetUserListQuery = { __typename?: 'Query' } & {
  getUserList?: Maybe<Array<Maybe<{ __typename?: 'AuthUser' } & LaForgeUserListFieldsFragment>>>;
};

export type LaForgeUpdateUserMutationVariables = Exact<{
  userId: Scalars['String'];
  firstName?: Maybe<Scalars['String']>;
  lastName?: Maybe<Scalars['String']>;
  email: Scalars['String'];
  phone?: Maybe<Scalars['String']>;
  company?: Maybe<Scalars['String']>;
  occupation?: Maybe<Scalars['String']>;
  role: LaForgeRoleLevel;
  provider: LaForgeProviderType;
}>;

export type LaForgeUpdateUserMutation = { __typename?: 'Mutation' } & {
  modifyAdminUserInfo?: Maybe<{ __typename?: 'AuthUser' } & LaForgeUserListFieldsFragment>;
};

export type LaForgeCreateUserMutationVariables = Exact<{
  username: Scalars['String'];
  password: Scalars['String'];
  provider: LaForgeProviderType;
  role: LaForgeRoleLevel;
}>;

export type LaForgeCreateUserMutation = { __typename?: 'Mutation' } & {
  createUser?: Maybe<{ __typename?: 'AuthUser' } & LaForgeUserListFieldsFragment>;
};

export type LaForgeGetAgentTasksQueryVariables = Exact<{
  proStepId: Scalars['String'];
}>;

export type LaForgeGetAgentTasksQuery = { __typename?: 'Query' } & {
  getAgentTasks?: Maybe<Array<Maybe<{ __typename?: 'AgentTask' } & LaForgeAgentTaskFieldsFragment>>>;
};

export type LaForgeGetCurrentUserQueryVariables = Exact<{ [key: string]: never }>;

export type LaForgeGetCurrentUserQuery = { __typename?: 'Query' } & {
  currentUser?: Maybe<{ __typename?: 'AuthUser' } & LaForgeAuthUserFieldsFragment>;
};

export type LaForgeGetBuildTreeQueryVariables = Exact<{
  buildId: Scalars['String'];
}>;

export type LaForgeGetBuildTreeQuery = { __typename?: 'Query' } & {
  build?: Maybe<
    { __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'> & {
        BuildToLatestBuildCommit?: Maybe<{ __typename?: 'BuildCommit' } & Pick<LaForgeBuildCommit, 'id'>>;
        buildToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
        buildToTeam: Array<
          Maybe<
            { __typename?: 'Team' } & Pick<LaForgeTeam, 'id' | 'team_number'> & {
                TeamToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                TeamToPlan: { __typename?: 'Plan' } & Pick<LaForgePlan, 'id'> & {
                    PlanToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                  };
                TeamToProvisionedNetwork: Array<
                  Maybe<
                    { __typename?: 'ProvisionedNetwork' } & Pick<LaForgeProvisionedNetwork, 'id' | 'name' | 'cidr'> & {
                        ProvisionedNetworkToNetwork: { __typename?: 'Network' } & Pick<LaForgeNetwork, 'id' | 'vdi_visible'> & {
                            vars?: Maybe<Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>>;
                            tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                          };
                        ProvisionedNetworkToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                        ProvisionedNetworkToPlan: { __typename?: 'Plan' } & Pick<LaForgePlan, 'id'> & {
                            PlanToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                          };
                        ProvisionedNetworkToProvisionedHost: Array<
                          Maybe<
                            { __typename?: 'ProvisionedHost' } & Pick<LaForgeProvisionedHost, 'id' | 'subnet_ip'> & {
                                ProvisionedHostToHost: { __typename?: 'Host' } & Pick<
                                  LaForgeHost,
                                  | 'id'
                                  | 'hostname'
                                  | 'description'
                                  | 'OS'
                                  | 'allow_mac_changes'
                                  | 'exposed_tcp_ports'
                                  | 'exposed_udp_ports'
                                  | 'user_groups'
                                  | 'override_password'
                                > & {
                                    vars?: Maybe<Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>>;
                                    tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                  };
                                ProvisionedHostToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                                ProvisionedHostToPlan: { __typename?: 'Plan' } & Pick<LaForgePlan, 'id'> & {
                                    PlanToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                                  };
                                ProvisionedHostToProvisioningStep: Array<
                                  Maybe<
                                    { __typename?: 'ProvisioningStep' } & Pick<LaForgeProvisioningStep, 'id' | 'type' | 'step_number'> & {
                                        ProvisioningStepToScript?: Maybe<
                                          { __typename?: 'Script' } & Pick<
                                            LaForgeScript,
                                            'id' | 'name' | 'language' | 'description' | 'source' | 'source_type' | 'disabled' | 'args'
                                          > & {
                                              vars?: Maybe<
                                                Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>
                                              >;
                                              tags?: Maybe<Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>>;
                                            }
                                        >;
                                        ProvisioningStepToCommand?: Maybe<
                                          { __typename?: 'Command' } & Pick<
                                            LaForgeCommand,
                                            'id' | 'name' | 'description' | 'program' | 'args' | 'disabled'
                                          > & {
                                              vars?: Maybe<
                                                Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>
                                              >;
                                              tags?: Maybe<Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>>;
                                            }
                                        >;
                                        ProvisioningStepToDNSRecord?: Maybe<
                                          { __typename?: 'DNSRecord' } & Pick<
                                            LaForgeDnsRecord,
                                            'id' | 'name' | 'values' | 'type' | 'zone' | 'disabled'
                                          > & {
                                              vars: Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>;
                                              tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                            }
                                        >;
                                        ProvisioningStepToFileDownload?: Maybe<
                                          { __typename?: 'FileDownload' } & Pick<
                                            LaForgeFileDownload,
                                            'id' | 'source' | 'sourceType' | 'destination' | 'disabled'
                                          > & { tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>> }
                                        >;
                                        ProvisioningStepToFileDelete?: Maybe<
                                          { __typename?: 'FileDelete' } & Pick<LaForgeFileDelete, 'id' | 'path'> & {
                                              tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                            }
                                        >;
                                        ProvisioningStepToFileExtract?: Maybe<
                                          { __typename?: 'FileExtract' } & Pick<
                                            LaForgeFileExtract,
                                            'id' | 'source' | 'destination' | 'type'
                                          > & { tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>> }
                                        >;
                                        ProvisioningStepToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                                        ProvisioningStepToPlan?: Maybe<
                                          { __typename?: 'Plan' } & Pick<LaForgePlan, 'id'> & {
                                              PlanToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                                            }
                                        >;
                                      }
                                  >
                                >;
                                ProvisionedHostToAgentStatus?: Maybe<{ __typename?: 'AgentStatus' } & Pick<LaForgeAgentStatus, 'clientId'>>;
                              }
                          >
                        >;
                      }
                  >
                >;
              }
          >
        >;
      }
  >;
};

export type LaForgeGetBuildPlansQueryVariables = Exact<{
  buildId: Scalars['String'];
}>;

export type LaForgeGetBuildPlansQuery = { __typename?: 'Query' } & {
  build?: Maybe<
    { __typename?: 'Build' } & Pick<LaForgeBuild, 'id'> & { buildToPlan: Array<Maybe<{ __typename?: 'Plan' } & LaForgePlanFieldsFragment>> }
  >;
};

export type LaForgeGetBuildCommitsQueryVariables = Exact<{
  buildId: Scalars['String'];
}>;

export type LaForgeGetBuildCommitsQuery = { __typename?: 'Query' } & {
  build?: Maybe<
    { __typename?: 'Build' } & Pick<LaForgeBuild, 'id'> & {
        BuildToBuildCommits: Array<Maybe<{ __typename?: 'BuildCommit' } & LaForgeBuildCommitFieldsFragment>>;
      }
  >;
};

export type LaForgeApproveBuildCommitMutationVariables = Exact<{
  buildCommitId: Scalars['String'];
}>;

export type LaForgeApproveBuildCommitMutation = { __typename?: 'Mutation' } & Pick<LaForgeMutation, 'approveCommit'>;

export type LaForgeCancelBuildCommitMutationVariables = Exact<{
  buildCommitId: Scalars['String'];
}>;

export type LaForgeCancelBuildCommitMutation = { __typename?: 'Mutation' } & Pick<LaForgeMutation, 'cancelCommit'>;

export type LaForgeGetEnvironmentQueryVariables = Exact<{
  envId: Scalars['String'];
}>;

export type LaForgeGetEnvironmentQuery = { __typename?: 'Query' } & {
  environment?: Maybe<
    { __typename?: 'Environment' } & Pick<
      LaForgeEnvironment,
      'id' | 'competition_id' | 'name' | 'description' | 'builder' | 'team_count' | 'revision' | 'admin_cidrs' | 'exposed_vdi_ports'
    > & {
        tags?: Maybe<Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>>;
        config?: Maybe<Array<Maybe<{ __typename?: 'configMap' } & Pick<LaForgeConfigMap, 'key' | 'value'>>>>;
        EnvironmentToUser: Array<Maybe<{ __typename?: 'User' } & Pick<LaForgeUser, 'id' | 'name' | 'uuid' | 'email'>>>;
        EnvironmentToRepository: Array<
          Maybe<{ __typename?: 'Repository' } & Pick<LaForgeRepository, 'id' | 'repo_url' | 'branch_name' | 'commit_info'>>
        >;
        EnvironmentToBuild: Array<
          Maybe<
            { __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'> & {
                buildToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
                buildToTeam: Array<
                  Maybe<
                    { __typename?: 'Team' } & Pick<LaForgeTeam, 'id' | 'team_number'> & {
                        TeamToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
                        TeamToProvisionedNetwork: Array<
                          Maybe<
                            { __typename?: 'ProvisionedNetwork' } & Pick<LaForgeProvisionedNetwork, 'id' | 'name' | 'cidr'> & {
                                ProvisionedNetworkToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
                                ProvisionedNetworkToNetwork: { __typename?: 'Network' } & Pick<LaForgeNetwork, 'id' | 'vdi_visible'> & {
                                    vars?: Maybe<Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>>;
                                    tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                  };
                                ProvisionedNetworkToProvisionedHost: Array<
                                  Maybe<
                                    { __typename?: 'ProvisionedHost' } & Pick<LaForgeProvisionedHost, 'id' | 'subnet_ip'> & {
                                        ProvisionedHostToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
                                        ProvisionedHostToHost: { __typename?: 'Host' } & Pick<
                                          LaForgeHost,
                                          | 'id'
                                          | 'hostname'
                                          | 'description'
                                          | 'OS'
                                          | 'allow_mac_changes'
                                          | 'exposed_tcp_ports'
                                          | 'exposed_udp_ports'
                                          | 'user_groups'
                                          | 'override_password'
                                        > & {
                                            vars?: Maybe<Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>>;
                                            tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                          };
                                        ProvisionedHostToProvisioningStep: Array<
                                          Maybe<
                                            { __typename?: 'ProvisioningStep' } & Pick<LaForgeProvisioningStep, 'id' | 'type'> & {
                                                ProvisioningStepToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
                                                ProvisioningStepToScript?: Maybe<
                                                  { __typename?: 'Script' } & Pick<
                                                    LaForgeScript,
                                                    | 'id'
                                                    | 'name'
                                                    | 'language'
                                                    | 'description'
                                                    | 'source'
                                                    | 'source_type'
                                                    | 'disabled'
                                                    | 'args'
                                                  > & {
                                                      vars?: Maybe<
                                                        Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>
                                                      >;
                                                      tags?: Maybe<
                                                        Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>
                                                      >;
                                                    }
                                                >;
                                                ProvisioningStepToCommand?: Maybe<
                                                  { __typename?: 'Command' } & Pick<
                                                    LaForgeCommand,
                                                    'id' | 'name' | 'description' | 'program' | 'args' | 'disabled'
                                                  > & {
                                                      vars?: Maybe<
                                                        Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>
                                                      >;
                                                      tags?: Maybe<
                                                        Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>
                                                      >;
                                                    }
                                                >;
                                                ProvisioningStepToDNSRecord?: Maybe<
                                                  { __typename?: 'DNSRecord' } & Pick<
                                                    LaForgeDnsRecord,
                                                    'id' | 'name' | 'values' | 'type' | 'zone' | 'disabled'
                                                  > & {
                                                      vars: Array<
                                                        Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>
                                                      >;
                                                      tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                                    }
                                                >;
                                                ProvisioningStepToFileDownload?: Maybe<
                                                  { __typename?: 'FileDownload' } & Pick<
                                                    LaForgeFileDownload,
                                                    'id' | 'source' | 'sourceType' | 'destination' | 'disabled'
                                                  > & {
                                                      tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                                    }
                                                >;
                                                ProvisioningStepToFileDelete?: Maybe<
                                                  { __typename?: 'FileDelete' } & Pick<LaForgeFileDelete, 'id' | 'path'> & {
                                                      tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                                    }
                                                >;
                                                ProvisioningStepToFileExtract?: Maybe<
                                                  { __typename?: 'FileExtract' } & Pick<
                                                    LaForgeFileExtract,
                                                    'id' | 'source' | 'destination' | 'type'
                                                  > & {
                                                      tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                                    }
                                                >;
                                              }
                                          >
                                        >;
                                      }
                                  >
                                >;
                              }
                          >
                        >;
                      }
                  >
                >;
              }
          >
        >;
      }
  >;
};

export type LaForgeGetEnvironmentsQueryVariables = Exact<{ [key: string]: never }>;

export type LaForgeGetEnvironmentsQuery = { __typename?: 'Query' } & {
  environments?: Maybe<
    Array<
      Maybe<
        { __typename?: 'Environment' } & Pick<LaForgeEnvironment, 'id' | 'name' | 'competition_id' | 'revision'> & {
            EnvironmentToBuild: Array<Maybe<{ __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'>>>;
          }
      >
    >
  >;
};

export type LaForgeGetEnvironmentInfoQueryVariables = Exact<{
  envId: Scalars['String'];
}>;

export type LaForgeGetEnvironmentInfoQuery = { __typename?: 'Query' } & {
  environment?: Maybe<
    { __typename?: 'Environment' } & Pick<
      LaForgeEnvironment,
      'id' | 'competition_id' | 'name' | 'description' | 'builder' | 'team_count' | 'revision' | 'admin_cidrs' | 'exposed_vdi_ports'
    > & {
        tags?: Maybe<Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>>;
        config?: Maybe<Array<Maybe<{ __typename?: 'configMap' } & Pick<LaForgeConfigMap, 'key' | 'value'>>>>;
        EnvironmentToUser: Array<Maybe<{ __typename?: 'User' } & Pick<LaForgeUser, 'id' | 'name' | 'uuid' | 'email'>>>;
        EnvironmentToBuild: Array<
          Maybe<
            { __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'> & {
                buildToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
              }
          >
        >;
      }
  >;
};

export type LaForgeStatusFieldsFragment = { __typename?: 'Status' } & Pick<
  LaForgeStatus,
  'id' | 'state' | 'started_at' | 'ended_at' | 'failed' | 'completed' | 'error'
>;

export type LaForgeAgentStatusFieldsFragment = { __typename?: 'AgentStatus' } & Pick<
  LaForgeAgentStatus,
  | 'clientId'
  | 'hostname'
  | 'upTime'
  | 'bootTime'
  | 'numProcs'
  | 'OS'
  | 'hostID'
  | 'load1'
  | 'load5'
  | 'load15'
  | 'totalMem'
  | 'freeMem'
  | 'usedMem'
  | 'timestamp'
>;

export type LaForgePlanFieldsFragment = { __typename?: 'Plan' } & Pick<LaForgePlan, 'id' | 'step_number' | 'type'> & {
    PlanToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
    PlanToPlanDiffs: Array<Maybe<{ __typename?: 'PlanDiff' } & LaForgePlanDiffFieldsFragment>>;
  };

export type LaForgePlanDiffFieldsFragment = { __typename?: 'PlanDiff' } & Pick<LaForgePlanDiff, 'id' | 'revision' | 'new_state'>;

export type LaForgeBuildCommitFieldsFragment = { __typename?: 'BuildCommit' } & Pick<
  LaForgeBuildCommit,
  'id' | 'revision' | 'type' | 'state'
> & { BuildCommitToPlanDiffs: Array<Maybe<{ __typename?: 'PlanDiff' } & LaForgePlanDiffFieldsFragment>> };

export type LaForgeAuthUserFieldsFragment = { __typename?: 'AuthUser' } & Pick<
  LaForgeAuthUser,
  'id' | 'username' | 'role' | 'provider' | 'first_name' | 'last_name' | 'email' | 'phone' | 'company' | 'occupation' | 'publicKey'
>;

export type LaForgeAgentTaskFieldsFragment = { __typename?: 'AgentTask' } & Pick<
  LaForgeAgentTask,
  'id' | 'state' | 'command' | 'args' | 'number' | 'output' | 'error_message'
>;

export type LaForgePageInfoFieldsFragment = { __typename?: 'LaForgePageInfo' } & Pick<LaForgeLaForgePageInfo, 'total' | 'nextOffset'>;

export type LaForgeUserListFieldsFragment = { __typename?: 'AuthUser' } & Pick<
  LaForgeAuthUser,
  'id' | 'first_name' | 'last_name' | 'username' | 'provider' | 'role' | 'email' | 'phone' | 'company' | 'occupation'
>;

export type LaForgeRebuildMutationVariables = Exact<{
  rootPlans: Array<Maybe<Scalars['String']>> | Maybe<Scalars['String']>;
}>;

export type LaForgeRebuildMutation = { __typename?: 'Mutation' } & Pick<LaForgeMutation, 'rebuild'>;

export type LaForgeDeleteBuildMutationVariables = Exact<{
  buildId: Scalars['String'];
}>;

export type LaForgeDeleteBuildMutation = { __typename?: 'Mutation' } & Pick<LaForgeMutation, 'deleteBuild'>;

export type LaForgeExecuteBuildMutationVariables = Exact<{
  buildId: Scalars['String'];
}>;

export type LaForgeExecuteBuildMutation = { __typename?: 'Mutation' } & {
  executePlan?: Maybe<{ __typename?: 'Build' } & Pick<LaForgeBuild, 'id'>>;
};

export type LaForgeCreateBuildMutationVariables = Exact<{
  envId: Scalars['String'];
}>;

export type LaForgeCreateBuildMutation = { __typename?: 'Mutation' } & {
  createBuild?: Maybe<{ __typename?: 'Build' } & Pick<LaForgeBuild, 'id'>>;
};

export type LaForgeModifyCurrentUserMutationVariables = Exact<{
  firstName?: Maybe<Scalars['String']>;
  lastName?: Maybe<Scalars['String']>;
  email?: Maybe<Scalars['String']>;
  phone?: Maybe<Scalars['String']>;
  company?: Maybe<Scalars['String']>;
  occupation?: Maybe<Scalars['String']>;
}>;

export type LaForgeModifyCurrentUserMutation = { __typename?: 'Mutation' } & {
  modifySelfUserInfo?: Maybe<{ __typename?: 'AuthUser' } & LaForgeAuthUserFieldsFragment>;
};

export type LaForgeCreateEnvironmentFromGitMutationVariables = Exact<{
  repoURL: Scalars['String'];
  repoName: Scalars['String'];
  branchName: Scalars['String'];
  envFilePath: Scalars['String'];
}>;

export type LaForgeCreateEnvironmentFromGitMutation = { __typename?: 'Mutation' } & {
  createEnviromentFromRepo: Array<Maybe<{ __typename?: 'Environment' } & Pick<LaForgeEnvironment, 'id'>>>;
};

export type LaForgeUpdateEnviromentViaPullMutationVariables = Exact<{
  envId: Scalars['String'];
}>;

export type LaForgeUpdateEnviromentViaPullMutation = { __typename?: 'Mutation' } & {
  updateEnviromentViaPull: Array<Maybe<{ __typename?: 'Environment' } & Pick<LaForgeEnvironment, 'id'>>>;
};

export type LaForgeGetStatusQueryVariables = Exact<{
  statusId: Scalars['String'];
}>;

export type LaForgeGetStatusQuery = { __typename?: 'Query' } & { status?: Maybe<{ __typename?: 'Status' } & LaForgeStatusFieldsFragment> };

export type LaForgeGetAgentStatusQueryVariables = Exact<{
  clientId: Scalars['String'];
}>;

export type LaForgeGetAgentStatusQuery = { __typename?: 'Query' } & {
  agentStatus?: Maybe<{ __typename?: 'AgentStatus' } & LaForgeAgentStatusFieldsFragment>;
};

export type LaForgeGetAllPlanStatusesQueryVariables = Exact<{
  buildId: Scalars['String'];
  count: Scalars['Int'];
  offset: Scalars['Int'];
}>;

export type LaForgeGetAllPlanStatusesQuery = { __typename?: 'Query' } & {
  getAllPlanStatus?: Maybe<
    { __typename?: 'StatusBatch' } & {
      statuses: Array<Maybe<{ __typename?: 'Status' } & LaForgeStatusFieldsFragment>>;
      pageInfo: { __typename?: 'LaForgePageInfo' } & LaForgePageInfoFieldsFragment;
    }
  >;
};

export type LaForgeGetAllAgentStatusesQueryVariables = Exact<{
  buildId: Scalars['String'];
  count: Scalars['Int'];
  offset: Scalars['Int'];
}>;

export type LaForgeGetAllAgentStatusesQuery = { __typename?: 'Query' } & {
  getAllAgentStatus?: Maybe<
    { __typename?: 'AgentStatusBatch' } & {
      agentStatuses: Array<Maybe<{ __typename?: 'AgentStatus' } & LaForgeAgentStatusFieldsFragment>>;
      pageInfo: { __typename?: 'LaForgePageInfo' } & LaForgePageInfoFieldsFragment;
    }
  >;
};

export type LaForgeSubscribeUpdatedStatusSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedStatusSubscription = { __typename?: 'Subscription' } & {
  updatedStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
};

export type LaForgeSubscribeUpdatedAgentStatusSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedAgentStatusSubscription = { __typename?: 'Subscription' } & {
  updatedAgentStatus: { __typename?: 'AgentStatus' } & LaForgeAgentStatusFieldsFragment;
};

export type LaForgeSubscribeUpdatedServerTaskSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedServerTaskSubscription = { __typename?: 'Subscription' } & {
  updatedServerTask: { __typename?: 'ServerTask' } & LaForgeServerTaskFieldsFragment;
};

export type LaForgeSubscribeUpdatedBuildSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedBuildSubscription = { __typename?: 'Subscription' } & {
  updatedBuild: { __typename?: 'Build' } & Pick<LaForgeBuild, 'id'> & {
      BuildToLatestBuildCommit?: Maybe<{ __typename?: 'BuildCommit' } & Pick<LaForgeBuildCommit, 'id'>>;
    };
};

export type LaForgeSubscribeUpdatedBuildCommitSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedBuildCommitSubscription = { __typename?: 'Subscription' } & {
  updatedCommit: { __typename?: 'BuildCommit' } & LaForgeBuildCommitFieldsFragment;
};

export type LaForgeSubscribeUpdatedAgentTaskSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedAgentTaskSubscription = { __typename?: 'Subscription' } & {
  updatedAgentTask: { __typename?: 'AgentTask' } & LaForgeAgentTaskFieldsFragment;
};

export type LaForgeServerTaskFieldsFragment = { __typename?: 'ServerTask' } & Pick<
  LaForgeServerTask,
  'id' | 'type' | 'start_time' | 'end_time' | 'errors' | 'log_file_path'
> & {
    ServerTaskToStatus: { __typename?: 'Status' } & LaForgeStatusFieldsFragment;
    ServerTaskToEnvironment?: Maybe<{ __typename?: 'Environment' } & Pick<LaForgeEnvironment, 'id' | 'name'>>;
    ServerTaskToBuild?: Maybe<{ __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'>>;
  };

export type LaForgeGetCurrentUserTasksQueryVariables = Exact<{ [key: string]: never }>;

export type LaForgeGetCurrentUserTasksQuery = { __typename?: 'Query' } & {
  getCurrentUserTasks?: Maybe<Array<Maybe<{ __typename?: 'ServerTask' } & LaForgeServerTaskFieldsFragment>>>;
};

export const AgentStatusFieldsFragmentDoc = gql`
  fragment AgentStatusFields on AgentStatus {
    clientId
    hostname
    upTime
    bootTime
    numProcs
    OS
    hostID
    load1
    load5
    load15
    totalMem
    freeMem
    usedMem
    timestamp
  }
`;
export const StatusFieldsFragmentDoc = gql`
  fragment StatusFields on Status {
    id
    state
    started_at
    ended_at
    failed
    completed
    error
  }
`;
export const PlanDiffFieldsFragmentDoc = gql`
  fragment PlanDiffFields on PlanDiff {
    id
    revision
    new_state
  }
`;
export const PlanFieldsFragmentDoc = gql`
  fragment PlanFields on Plan {
    id
    step_number
    type
    PlanToStatus {
      ...StatusFields
    }
    PlanToPlanDiffs {
      ...PlanDiffFields
    }
  }
  ${StatusFieldsFragmentDoc}
  ${PlanDiffFieldsFragmentDoc}
`;
export const BuildCommitFieldsFragmentDoc = gql`
  fragment BuildCommitFields on BuildCommit {
    id
    revision
    type
    state
    BuildCommitToPlanDiffs {
      ...PlanDiffFields
    }
  }
  ${PlanDiffFieldsFragmentDoc}
`;
export const AuthUserFieldsFragmentDoc = gql`
  fragment AuthUserFields on AuthUser {
    id
    username
    role
    provider
    first_name
    last_name
    email
    phone
    company
    occupation
    publicKey
  }
`;
export const AgentTaskFieldsFragmentDoc = gql`
  fragment AgentTaskFields on AgentTask {
    id
    state
    command
    args
    number
    output
    error_message
  }
`;
export const PageInfoFieldsFragmentDoc = gql`
  fragment PageInfoFields on LaForgePageInfo {
    total
    nextOffset
  }
`;
export const UserListFieldsFragmentDoc = gql`
  fragment UserListFields on AuthUser {
    id
    first_name
    last_name
    username
    provider
    role
    email
    phone
    company
    occupation
  }
`;
export const ServerTaskFieldsFragmentDoc = gql`
  fragment ServerTaskFields on ServerTask {
    id
    type
    start_time
    end_time
    errors
    log_file_path
    ServerTaskToStatus {
      ...StatusFields
    }
    ServerTaskToEnvironment {
      id
      name
    }
    ServerTaskToBuild {
      id
      revision
    }
  }
  ${StatusFieldsFragmentDoc}
`;
export const GetUserListDocument = gql`
  query GetUserList {
    getUserList {
      ...UserListFields
    }
  }
  ${UserListFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetUserListGQL extends Apollo.Query<LaForgeGetUserListQuery, LaForgeGetUserListQueryVariables> {
  document = GetUserListDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const UpdateUserDocument = gql`
  mutation UpdateUser(
    $userId: String!
    $firstName: String
    $lastName: String
    $email: String!
    $phone: String
    $company: String
    $occupation: String
    $role: RoleLevel!
    $provider: ProviderType!
  ) {
    modifyAdminUserInfo(
      userID: $userId
      firstName: $firstName
      lastName: $lastName
      email: $email
      phone: $phone
      company: $company
      occupation: $occupation
      role: $role
      provider: $provider
    ) {
      ...UserListFields
    }
  }
  ${UserListFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeUpdateUserGQL extends Apollo.Mutation<LaForgeUpdateUserMutation, LaForgeUpdateUserMutationVariables> {
  document = UpdateUserDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const CreateUserDocument = gql`
  mutation CreateUser($username: String!, $password: String!, $provider: ProviderType!, $role: RoleLevel!) {
    createUser(username: $username, password: $password, provider: $provider, role: $role) {
      ...UserListFields
    }
  }
  ${UserListFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeCreateUserGQL extends Apollo.Mutation<LaForgeCreateUserMutation, LaForgeCreateUserMutationVariables> {
  document = CreateUserDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetAgentTasksDocument = gql`
  query GetAgentTasks($proStepId: String!) {
    getAgentTasks(proStepUUID: $proStepId) {
      ...AgentTaskFields
    }
  }
  ${AgentTaskFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetAgentTasksGQL extends Apollo.Query<LaForgeGetAgentTasksQuery, LaForgeGetAgentTasksQueryVariables> {
  document = GetAgentTasksDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetCurrentUserDocument = gql`
  query GetCurrentUser {
    currentUser {
      ...AuthUserFields
    }
  }
  ${AuthUserFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetCurrentUserGQL extends Apollo.Query<LaForgeGetCurrentUserQuery, LaForgeGetCurrentUserQueryVariables> {
  document = GetCurrentUserDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetBuildTreeDocument = gql`
  query GetBuildTree($buildId: String!) {
    build(buildUUID: $buildId) {
      id
      revision
      BuildToLatestBuildCommit {
        id
      }
      buildToStatus {
        ...StatusFields
      }
      buildToTeam {
        id
        team_number
        TeamToStatus {
          id
        }
        TeamToPlan {
          id
          PlanToStatus {
            id
          }
        }
        TeamToProvisionedNetwork {
          id
          name
          cidr
          ProvisionedNetworkToNetwork {
            id
            vdi_visible
            vars {
              key
              value
            }
            tags {
              key
              value
            }
          }
          ProvisionedNetworkToStatus {
            id
          }
          ProvisionedNetworkToPlan {
            id
            PlanToStatus {
              id
            }
          }
          ProvisionedNetworkToProvisionedHost {
            id
            subnet_ip
            ProvisionedHostToHost {
              id
              hostname
              description
              OS
              allow_mac_changes
              exposed_tcp_ports
              exposed_udp_ports
              user_groups
              override_password
              vars {
                key
                value
              }
              tags {
                key
                value
              }
            }
            ProvisionedHostToStatus {
              id
            }
            ProvisionedHostToPlan {
              id
              PlanToStatus {
                id
              }
            }
            ProvisionedHostToProvisioningStep {
              id
              type
              step_number
              ProvisioningStepToScript {
                id
                name
                language
                description
                source
                source_type
                disabled
                args
                vars {
                  key
                  value
                }
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToCommand {
                id
                name
                description
                program
                args
                disabled
                vars {
                  key
                  value
                }
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToDNSRecord {
                id
                name
                values
                type
                zone
                disabled
                vars {
                  key
                  value
                }
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToFileDownload {
                id
                source
                sourceType
                destination
                disabled
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToFileDelete {
                id
                path
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToFileExtract {
                id
                source
                destination
                type
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToStatus {
                id
              }
              ProvisioningStepToPlan {
                id
                PlanToStatus {
                  id
                }
              }
            }
            ProvisionedHostToAgentStatus {
              clientId
            }
          }
        }
      }
    }
  }
  ${StatusFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetBuildTreeGQL extends Apollo.Query<LaForgeGetBuildTreeQuery, LaForgeGetBuildTreeQueryVariables> {
  document = GetBuildTreeDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetBuildPlansDocument = gql`
  query GetBuildPlans($buildId: String!) {
    build(buildUUID: $buildId) {
      id
      buildToPlan {
        ...PlanFields
      }
    }
  }
  ${PlanFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetBuildPlansGQL extends Apollo.Query<LaForgeGetBuildPlansQuery, LaForgeGetBuildPlansQueryVariables> {
  document = GetBuildPlansDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetBuildCommitsDocument = gql`
  query GetBuildCommits($buildId: String!) {
    build(buildUUID: $buildId) {
      id
      BuildToBuildCommits {
        ...BuildCommitFields
      }
    }
  }
  ${BuildCommitFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetBuildCommitsGQL extends Apollo.Query<LaForgeGetBuildCommitsQuery, LaForgeGetBuildCommitsQueryVariables> {
  document = GetBuildCommitsDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const ApproveBuildCommitDocument = gql`
  mutation ApproveBuildCommit($buildCommitId: String!) {
    approveCommit(commitUUID: $buildCommitId)
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeApproveBuildCommitGQL extends Apollo.Mutation<
  LaForgeApproveBuildCommitMutation,
  LaForgeApproveBuildCommitMutationVariables
> {
  document = ApproveBuildCommitDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const CancelBuildCommitDocument = gql`
  mutation CancelBuildCommit($buildCommitId: String!) {
    cancelCommit(commitUUID: $buildCommitId)
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeCancelBuildCommitGQL extends Apollo.Mutation<
  LaForgeCancelBuildCommitMutation,
  LaForgeCancelBuildCommitMutationVariables
> {
  document = CancelBuildCommitDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetEnvironmentDocument = gql`
  query GetEnvironment($envId: String!) {
    environment(envUUID: $envId) {
      id
      competition_id
      name
      description
      builder
      team_count
      revision
      admin_cidrs
      exposed_vdi_ports
      tags {
        key
        value
      }
      config {
        key
        value
      }
      EnvironmentToUser {
        id
        name
        uuid
        email
      }
      EnvironmentToRepository {
        id
        repo_url
        branch_name
        commit_info
      }
      EnvironmentToBuild {
        id
        revision
        buildToStatus {
          ...StatusFields
        }
        buildToTeam {
          id
          team_number
          TeamToStatus {
            ...StatusFields
          }
          TeamToProvisionedNetwork {
            id
            name
            cidr
            ProvisionedNetworkToStatus {
              ...StatusFields
            }
            ProvisionedNetworkToNetwork {
              id
              vdi_visible
              vars {
                key
                value
              }
              tags {
                key
                value
              }
            }
            ProvisionedNetworkToProvisionedHost {
              id
              subnet_ip
              ProvisionedHostToStatus {
                ...StatusFields
              }
              ProvisionedHostToHost {
                id
                hostname
                description
                OS
                allow_mac_changes
                exposed_tcp_ports
                exposed_udp_ports
                user_groups
                override_password
                vars {
                  key
                  value
                }
                tags {
                  key
                  value
                }
              }
              ProvisionedHostToProvisioningStep {
                id
                type
                ProvisioningStepToStatus {
                  ...StatusFields
                }
                ProvisioningStepToScript {
                  id
                  name
                  language
                  description
                  source
                  source_type
                  disabled
                  args
                  vars {
                    key
                    value
                  }
                  tags {
                    key
                    value
                  }
                }
                ProvisioningStepToCommand {
                  id
                  name
                  description
                  program
                  args
                  disabled
                  vars {
                    key
                    value
                  }
                  tags {
                    key
                    value
                  }
                }
                ProvisioningStepToDNSRecord {
                  id
                  name
                  values
                  type
                  zone
                  disabled
                  vars {
                    key
                    value
                  }
                  tags {
                    key
                    value
                  }
                }
                ProvisioningStepToFileDownload {
                  id
                  source
                  sourceType
                  destination
                  disabled
                  tags {
                    key
                    value
                  }
                }
                ProvisioningStepToFileDelete {
                  id
                  path
                  tags {
                    key
                    value
                  }
                }
                ProvisioningStepToFileExtract {
                  id
                  source
                  destination
                  type
                  tags {
                    key
                    value
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  ${StatusFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetEnvironmentGQL extends Apollo.Query<LaForgeGetEnvironmentQuery, LaForgeGetEnvironmentQueryVariables> {
  document = GetEnvironmentDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetEnvironmentsDocument = gql`
  query GetEnvironments {
    environments {
      id
      name
      competition_id
      revision
      EnvironmentToBuild {
        id
        revision
      }
    }
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetEnvironmentsGQL extends Apollo.Query<LaForgeGetEnvironmentsQuery, LaForgeGetEnvironmentsQueryVariables> {
  document = GetEnvironmentsDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetEnvironmentInfoDocument = gql`
  query GetEnvironmentInfo($envId: String!) {
    environment(envUUID: $envId) {
      id
      competition_id
      name
      description
      builder
      team_count
      revision
      admin_cidrs
      exposed_vdi_ports
      tags {
        key
        value
      }
      config {
        key
        value
      }
      EnvironmentToUser {
        id
        name
        uuid
        email
      }
      EnvironmentToBuild {
        id
        revision
        buildToStatus {
          ...StatusFields
        }
      }
    }
  }
  ${StatusFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetEnvironmentInfoGQL extends Apollo.Query<LaForgeGetEnvironmentInfoQuery, LaForgeGetEnvironmentInfoQueryVariables> {
  document = GetEnvironmentInfoDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const RebuildDocument = gql`
  mutation Rebuild($rootPlans: [String]!) {
    rebuild(rootPlans: $rootPlans)
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeRebuildGQL extends Apollo.Mutation<LaForgeRebuildMutation, LaForgeRebuildMutationVariables> {
  document = RebuildDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const DeleteBuildDocument = gql`
  mutation DeleteBuild($buildId: String!) {
    deleteBuild(buildUUID: $buildId)
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeDeleteBuildGQL extends Apollo.Mutation<LaForgeDeleteBuildMutation, LaForgeDeleteBuildMutationVariables> {
  document = DeleteBuildDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const ExecuteBuildDocument = gql`
  mutation ExecuteBuild($buildId: String!) {
    executePlan(buildUUID: $buildId) {
      id
    }
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeExecuteBuildGQL extends Apollo.Mutation<LaForgeExecuteBuildMutation, LaForgeExecuteBuildMutationVariables> {
  document = ExecuteBuildDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const CreateBuildDocument = gql`
  mutation CreateBuild($envId: String!) {
    createBuild(envUUID: $envId, renderFiles: true) {
      id
    }
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeCreateBuildGQL extends Apollo.Mutation<LaForgeCreateBuildMutation, LaForgeCreateBuildMutationVariables> {
  document = CreateBuildDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const ModifyCurrentUserDocument = gql`
  mutation ModifyCurrentUser($firstName: String, $lastName: String, $email: String, $phone: String, $company: String, $occupation: String) {
    modifySelfUserInfo(
      firstName: $firstName
      lastName: $lastName
      email: $email
      phone: $phone
      company: $company
      occupation: $occupation
    ) {
      ...AuthUserFields
    }
  }
  ${AuthUserFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeModifyCurrentUserGQL extends Apollo.Mutation<
  LaForgeModifyCurrentUserMutation,
  LaForgeModifyCurrentUserMutationVariables
> {
  document = ModifyCurrentUserDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const CreateEnvironmentFromGitDocument = gql`
  mutation CreateEnvironmentFromGit($repoURL: String!, $repoName: String!, $branchName: String!, $envFilePath: String!) {
    createEnviromentFromRepo(repoURL: $repoURL, repoName: $repoName, branchName: $branchName, envFilePath: $envFilePath) {
      id
    }
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeCreateEnvironmentFromGitGQL extends Apollo.Mutation<
  LaForgeCreateEnvironmentFromGitMutation,
  LaForgeCreateEnvironmentFromGitMutationVariables
> {
  document = CreateEnvironmentFromGitDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const UpdateEnviromentViaPullDocument = gql`
  mutation UpdateEnviromentViaPull($envId: String!) {
    updateEnviromentViaPull(envUUID: $envId) {
      id
    }
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeUpdateEnviromentViaPullGQL extends Apollo.Mutation<
  LaForgeUpdateEnviromentViaPullMutation,
  LaForgeUpdateEnviromentViaPullMutationVariables
> {
  document = UpdateEnviromentViaPullDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetStatusDocument = gql`
  query GetStatus($statusId: String!) {
    status(statusUUID: $statusId) {
      ...StatusFields
    }
  }
  ${StatusFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetStatusGQL extends Apollo.Query<LaForgeGetStatusQuery, LaForgeGetStatusQueryVariables> {
  document = GetStatusDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetAgentStatusDocument = gql`
  query GetAgentStatus($clientId: String!) {
    agentStatus(clientId: $clientId) {
      ...AgentStatusFields
    }
  }
  ${AgentStatusFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetAgentStatusGQL extends Apollo.Query<LaForgeGetAgentStatusQuery, LaForgeGetAgentStatusQueryVariables> {
  document = GetAgentStatusDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetAllPlanStatusesDocument = gql`
  query GetAllPlanStatuses($buildId: String!, $count: Int!, $offset: Int!) {
    getAllPlanStatus(buildUUID: $buildId, count: $count, offset: $offset) {
      statuses {
        ...StatusFields
      }
      pageInfo {
        ...PageInfoFields
      }
    }
  }
  ${StatusFieldsFragmentDoc}
  ${PageInfoFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetAllPlanStatusesGQL extends Apollo.Query<LaForgeGetAllPlanStatusesQuery, LaForgeGetAllPlanStatusesQueryVariables> {
  document = GetAllPlanStatusesDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetAllAgentStatusesDocument = gql`
  query GetAllAgentStatuses($buildId: String!, $count: Int!, $offset: Int!) {
    getAllAgentStatus(buildUUID: $buildId, count: $count, offset: $offset) {
      agentStatuses {
        ...AgentStatusFields
      }
      pageInfo {
        ...PageInfoFields
      }
    }
  }
  ${AgentStatusFieldsFragmentDoc}
  ${PageInfoFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetAllAgentStatusesGQL extends Apollo.Query<LaForgeGetAllAgentStatusesQuery, LaForgeGetAllAgentStatusesQueryVariables> {
  document = GetAllAgentStatusesDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const SubscribeUpdatedStatusDocument = gql`
  subscription SubscribeUpdatedStatus {
    updatedStatus {
      ...StatusFields
    }
  }
  ${StatusFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeSubscribeUpdatedStatusGQL extends Apollo.Subscription<
  LaForgeSubscribeUpdatedStatusSubscription,
  LaForgeSubscribeUpdatedStatusSubscriptionVariables
> {
  document = SubscribeUpdatedStatusDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const SubscribeUpdatedAgentStatusDocument = gql`
  subscription SubscribeUpdatedAgentStatus {
    updatedAgentStatus {
      ...AgentStatusFields
    }
  }
  ${AgentStatusFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeSubscribeUpdatedAgentStatusGQL extends Apollo.Subscription<
  LaForgeSubscribeUpdatedAgentStatusSubscription,
  LaForgeSubscribeUpdatedAgentStatusSubscriptionVariables
> {
  document = SubscribeUpdatedAgentStatusDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const SubscribeUpdatedServerTaskDocument = gql`
  subscription SubscribeUpdatedServerTask {
    updatedServerTask {
      ...ServerTaskFields
    }
  }
  ${ServerTaskFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeSubscribeUpdatedServerTaskGQL extends Apollo.Subscription<
  LaForgeSubscribeUpdatedServerTaskSubscription,
  LaForgeSubscribeUpdatedServerTaskSubscriptionVariables
> {
  document = SubscribeUpdatedServerTaskDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const SubscribeUpdatedBuildDocument = gql`
  subscription SubscribeUpdatedBuild {
    updatedBuild {
      id
      BuildToLatestBuildCommit {
        id
      }
    }
  }
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeSubscribeUpdatedBuildGQL extends Apollo.Subscription<
  LaForgeSubscribeUpdatedBuildSubscription,
  LaForgeSubscribeUpdatedBuildSubscriptionVariables
> {
  document = SubscribeUpdatedBuildDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const SubscribeUpdatedBuildCommitDocument = gql`
  subscription SubscribeUpdatedBuildCommit {
    updatedCommit {
      ...BuildCommitFields
    }
  }
  ${BuildCommitFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeSubscribeUpdatedBuildCommitGQL extends Apollo.Subscription<
  LaForgeSubscribeUpdatedBuildCommitSubscription,
  LaForgeSubscribeUpdatedBuildCommitSubscriptionVariables
> {
  document = SubscribeUpdatedBuildCommitDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const SubscribeUpdatedAgentTaskDocument = gql`
  subscription SubscribeUpdatedAgentTask {
    updatedAgentTask {
      ...AgentTaskFields
    }
  }
  ${AgentTaskFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeSubscribeUpdatedAgentTaskGQL extends Apollo.Subscription<
  LaForgeSubscribeUpdatedAgentTaskSubscription,
  LaForgeSubscribeUpdatedAgentTaskSubscriptionVariables
> {
  document = SubscribeUpdatedAgentTaskDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
export const GetCurrentUserTasksDocument = gql`
  query GetCurrentUserTasks {
    getCurrentUserTasks {
      ...ServerTaskFields
    }
  }
  ${ServerTaskFieldsFragmentDoc}
`;

@Injectable({
  providedIn: 'root'
})
export class LaForgeGetCurrentUserTasksGQL extends Apollo.Query<LaForgeGetCurrentUserTasksQuery, LaForgeGetCurrentUserTasksQueryVariables> {
  document = GetCurrentUserTasksDocument;

  constructor(apollo: Apollo.Apollo) {
    super(apollo);
  }
}
