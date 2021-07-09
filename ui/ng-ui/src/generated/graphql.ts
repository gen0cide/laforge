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
  Validate = 'VALIDATE'
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

export type LaForgeAuthUser = {
  __typename?: 'AuthUser';
  id: Scalars['ID'];
  username: Scalars['String'];
  password: Scalars['String'];
  role: LaForgeRoleLevel;
  provider: LaForgeProviderType;
};

export type LaForgeBuild = {
  __typename?: 'Build';
  id: Scalars['ID'];
  revision: Scalars['Int'];
  completed_plan: Scalars['Boolean'];
  buildToStatus: LaForgeStatus;
  buildToEnvironment: LaForgeEnvironment;
  buildToCompetition: LaForgeCompetition;
  buildToProvisionedNetwork: Array<Maybe<LaForgeProvisionedNetwork>>;
  buildToTeam: Array<Maybe<LaForgeTeam>>;
  buildToPlan: Array<Maybe<LaForgePlan>>;
};

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

export type LaForgeMutation = {
  __typename?: 'Mutation';
  loadEnviroment?: Maybe<Array<Maybe<LaForgeEnvironment>>>;
  createBuild?: Maybe<LaForgeBuild>;
  createUser: LaForgeAuthUser;
  deleteUser: Scalars['Boolean'];
  executePlan?: Maybe<LaForgeBuild>;
  deleteBuild: Scalars['Boolean'];
  createTask: Scalars['Boolean'];
  rebuild: Scalars['Boolean'];
};

export type LaForgeMutationLoadEnviromentArgs = {
  envFilePath: Scalars['String'];
};

export type LaForgeMutationCreateBuildArgs = {
  envUUID: Scalars['String'];
  renderFiles?: Scalars['Boolean'];
};

export type LaForgeMutationCreateUserArgs = {
  username: Scalars['String'];
  password: Scalars['String'];
  role: LaForgeRoleLevel;
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
  Deleted = 'DELETED'
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
  combined_output?: Maybe<Scalars['String']>;
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

export type LaForgeSubscription = {
  __typename?: 'Subscription';
  updatedAgentStatus: LaForgeAgentStatus;
  updatedStatus: LaForgeStatus;
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

export type LaForgeGetBuildTreeQueryVariables = Exact<{
  buildId: Scalars['String'];
}>;

export type LaForgeGetBuildTreeQuery = { __typename?: 'Query' } & {
  build?: Maybe<
    { __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'> & {
        buildToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
        buildToTeam: Array<
          Maybe<
            { __typename?: 'Team' } & Pick<LaForgeTeam, 'id' | 'team_number'> & {
                TeamToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                TeamToPlan: { __typename?: 'Plan' } & Pick<LaForgePlan, 'id'>;
                TeamToProvisionedNetwork: Array<
                  Maybe<
                    { __typename?: 'ProvisionedNetwork' } & Pick<LaForgeProvisionedNetwork, 'id' | 'name' | 'cidr'> & {
                        ProvisionedNetworkToNetwork: { __typename?: 'Network' } & Pick<LaForgeNetwork, 'id' | 'vdi_visible'> & {
                            vars?: Maybe<Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>>;
                            tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                          };
                        ProvisionedNetworkToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'id'>;
                        ProvisionedNetworkToPlan: { __typename?: 'Plan' } & Pick<LaForgePlan, 'id'>;
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
                                ProvisionedHostToPlan: { __typename?: 'Plan' } & Pick<LaForgePlan, 'id'>;
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
                                        ProvisioningStepToPlan?: Maybe<{ __typename?: 'Plan' } & Pick<LaForgePlan, 'id'>>;
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

export type LaForgeGetEnvironmentQueryVariables = Exact<{
  envId: Scalars['String'];
}>;

export type LaForgeGetEnvironmentQuery = { __typename?: 'Query' } & {
  environment?: Maybe<
    { __typename?: 'Environment' } & Pick<
      LaForgeEnvironment,
      'id' | 'competition_id' | 'name' | 'description' | 'builder' | 'team_count' | 'admin_cidrs' | 'exposed_vdi_ports'
    > & {
        tags?: Maybe<Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>>;
        config?: Maybe<Array<Maybe<{ __typename?: 'configMap' } & Pick<LaForgeConfigMap, 'key' | 'value'>>>>;
        EnvironmentToUser: Array<Maybe<{ __typename?: 'User' } & Pick<LaForgeUser, 'id' | 'name' | 'uuid' | 'email'>>>;
        EnvironmentToBuild: Array<
          Maybe<
            { __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'> & {
                buildToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'completed' | 'failed' | 'state'>;
                buildToTeam: Array<
                  Maybe<
                    { __typename?: 'Team' } & Pick<LaForgeTeam, 'id' | 'team_number'> & {
                        TeamToStatus: { __typename?: 'Status' } & Pick<LaForgeStatus, 'completed' | 'failed' | 'state'>;
                        TeamToProvisionedNetwork: Array<
                          Maybe<
                            { __typename?: 'ProvisionedNetwork' } & Pick<LaForgeProvisionedNetwork, 'id' | 'name' | 'cidr'> & {
                                ProvisionedNetworkToStatus: { __typename?: 'Status' } & Pick<
                                  LaForgeStatus,
                                  'completed' | 'failed' | 'state'
                                >;
                                ProvisionedNetworkToNetwork: { __typename?: 'Network' } & Pick<LaForgeNetwork, 'id' | 'vdi_visible'> & {
                                    vars?: Maybe<Array<Maybe<{ __typename?: 'varsMap' } & Pick<LaForgeVarsMap, 'key' | 'value'>>>>;
                                    tags: Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>;
                                  };
                                ProvisionedNetworkToProvisionedHost: Array<
                                  Maybe<
                                    { __typename?: 'ProvisionedHost' } & Pick<LaForgeProvisionedHost, 'id' | 'subnet_ip'> & {
                                        ProvisionedHostToStatus: { __typename?: 'Status' } & Pick<
                                          LaForgeStatus,
                                          'completed' | 'failed' | 'state'
                                        >;
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
                                                ProvisioningStepToStatus: { __typename?: 'Status' } & Pick<
                                                  LaForgeStatus,
                                                  'completed' | 'failed' | 'state'
                                                >;
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
        { __typename?: 'Environment' } & Pick<LaForgeEnvironment, 'id' | 'name' | 'competition_id'> & {
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
      'id' | 'competition_id' | 'name' | 'description' | 'builder' | 'team_count' | 'admin_cidrs' | 'exposed_vdi_ports'
    > & {
        tags?: Maybe<Array<Maybe<{ __typename?: 'tagMap' } & Pick<LaForgeTagMap, 'key' | 'value'>>>>;
        config?: Maybe<Array<Maybe<{ __typename?: 'configMap' } & Pick<LaForgeConfigMap, 'key' | 'value'>>>>;
        EnvironmentToUser: Array<Maybe<{ __typename?: 'User' } & Pick<LaForgeUser, 'id' | 'name' | 'uuid' | 'email'>>>;
        EnvironmentToBuild: Array<Maybe<{ __typename?: 'Build' } & Pick<LaForgeBuild, 'id' | 'revision'>>>;
      }
  >;
};

export type LaForgeRebuildMutationVariables = Exact<{
  rootPlans: Array<Maybe<Scalars['String']>> | Maybe<Scalars['String']>;
}>;

export type LaForgeRebuildMutation = { __typename?: 'Mutation' } & Pick<LaForgeMutation, 'rebuild'>;

export type LaForgeSubscribeUpdatedStatusSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedStatusSubscription = { __typename?: 'Subscription' } & {
  updatedStatus: { __typename?: 'Status' } & Pick<
    LaForgeStatus,
    'id' | 'state' | 'started_at' | 'ended_at' | 'failed' | 'completed' | 'error'
  >;
};

export type LaForgeSubscribeUpdatedAgentStatusSubscriptionVariables = Exact<{ [key: string]: never }>;

export type LaForgeSubscribeUpdatedAgentStatusSubscription = { __typename?: 'Subscription' } & {
  updatedAgentStatus: { __typename?: 'AgentStatus' } & Pick<
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
};

export const GetBuildTreeDocument = gql`
  query GetBuildTree($buildId: String!) {
    build(buildUUID: $buildId) {
      id
      revision
      buildToStatus {
        id
      }
      buildToTeam {
        id
        team_number
        TeamToStatus {
          id
        }
        TeamToPlan {
          id
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
export const GetEnvironmentDocument = gql`
  query GetEnvironment($envId: String!) {
    environment(envUUID: $envId) {
      id
      competition_id
      name
      description
      builder
      team_count
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
          completed
          failed
          state
        }
        buildToTeam {
          id
          team_number
          TeamToStatus {
            completed
            failed
            state
          }
          TeamToProvisionedNetwork {
            id
            name
            cidr
            ProvisionedNetworkToStatus {
              completed
              failed
              state
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
                completed
                failed
                state
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
                  completed
                  failed
                  state
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
      }
    }
  }
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
export const SubscribeUpdatedStatusDocument = gql`
  subscription SubscribeUpdatedStatus {
    updatedStatus {
      id
      state
      started_at
      ended_at
      failed
      completed
      error
    }
  }
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
  }
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
