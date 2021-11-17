import { Command } from './command.model';
import { ID, Plan, Status } from './common.model';
import { DNSRecord } from './dns.model';
import { FileDelete, FileDownload, FileExtract } from './file.model';
import { ProvisionedHost } from './host.model';
import { Script } from './script.model';

export enum ProvisioningStepType {
  Script,
  Command,
  DNSRecord,
  FileDelete,
  FileDownload,
  FileExtract,
  Undefined
}

export interface ProvisioningStep {
  id: ID;
  type: ProvisioningStepType;
  step_number: number;
  ProvisioningStepToStatus: Status;
  ProvisioningStepToProvisionedHost?: ProvisionedHost;
  ProvisioningStepToScript?: Script;
  ProvisioningStepToCommand?: Command;
  ProvisioningStepToDNSRecord?: DNSRecord;
  ProvisioningStepToFileDelete?: FileDelete;
  ProvisioningStepToFileDownload?: FileDownload;
  ProvisioningStepToFileExtract?: FileExtract;
  ProvisioningStepToPlan?: Plan;
}
