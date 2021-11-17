import { Build } from './environment.model';
import { ProvisionedHost } from './host.model';
import { ProvisionedNetwork } from './network.model';
import { ProvisioningStep } from './step.model';

export type ID = string | number;

export interface varsMap {
  key: string;
  value: string;
}

export interface configMap {
  key: string;
  value: string;
}

export interface tagMap {
  key: string;
  value: string;
}

export interface User {
  id?: ID;
  name: string;
  uuid?: string;
  email: string;
}

// export interface Tag {
//   id?: ID;
//   name: string;
//   description?: string;
// }

export interface Team {
  id: ID;
  team_number: number;
  TeamToBuild?: Build;
  TeamToStatus?: Status;
  TeamToProvisionedNetwork?: ProvisionedNetwork[];
  TeamToPlan?: Plan;
}

export enum ProvisionStatus {
  PLANNING,
  AWAITING,
  INPROGRESS,
  FAILED,
  COMPLETE,
  TAINTED,
  UNDEFINED,
  TODELETE,
  DELETEINPROGRESS,
  DELETED
}

export enum ProvisionStatusFor {
  Build,
  Team,
  Plan,
  ProvisionedNetwork,
  ProvisionedHost,
  ProvisioningStep,
  Undefined
}

export interface Status {
  state: ProvisionStatus;
  startedAt: string;
  endedAt: string;
  failed: boolean;
  completed: boolean;
  error: string;
  status_for?: ProvisionStatusFor;
}

export enum PlanType {
  start_build,
  start_team,
  provision_network,
  provision_host,
  execute_step,
  undefined
}

export interface Plan {
  id: ID;
  step_number: number;
  type: PlanType;
  PlanToStatus: Status;
  build_id?: ID;
  NextPlan?: Plan[];
  PrevPlan?: Plan[];
  PlanToBuild?: Build;
  PlanToTeam?: Team;
  PlanToProvisionedNetwork?: ProvisionedNetwork;
  PlanToProvisionedHost?: ProvisionedHost;
  PlanToProvisioningStep?: ProvisioningStep;
}
