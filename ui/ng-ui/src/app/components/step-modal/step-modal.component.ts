import { ChangeDetectorRef, Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import {
  LaForgeGetAgentTasksQuery,
  LaForgeProvisioningStep,
  LaForgeProvisioningStepType,
  LaForgeProvisionStatus,
  LaForgeStatus
} from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';
import { BehaviorSubject } from 'rxjs';

import { LaForgeGetAgentTasksGQL } from '../../../generated/graphql';

@Component({
  selector: 'app-network-modal',
  templateUrl: './step-modal.component.html',
  styleUrls: ['./step-modal.component.scss']
})
export class StepModalComponent implements OnInit {
  taskColumns: string[] = ['args', 'state'];
  failedChildren = false;
  agentTasks: BehaviorSubject<LaForgeGetAgentTasksQuery['getAgentTasks']>;

  constructor(
    public dialogRef: MatDialogRef<StepModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { provisioningStep: LaForgeProvisioningStep; planStatus: LaForgeStatus },
    private getAgentTasks: LaForgeGetAgentTasksGQL,
    private cdRef: ChangeDetectorRef,
    private envService: EnvironmentService
  ) {
    this.agentTasks = new BehaviorSubject([]);
  }

  ngOnInit(): void {
    this.getAgentTasks
      .fetch({
        proStepId: this.data.provisioningStep.id
      })
      .toPromise()
      .then(({ data, error, errors }) => {
        if (error) {
          return this.agentTasks.error(error);
        } else if (errors) {
          return this.agentTasks.error(errors);
        }
        const tasks = [...data.getAgentTasks];
        for (let i = 0; i < tasks.length; i++) {
          const updatedTask = this.envService.getAgentTask(tasks[i].id);
          if (updatedTask) tasks[i] = { ...updatedTask };
        }
        this.agentTasks.next(tasks.sort((a, b) => a.number - b.number));
      }, this.agentTasks.error);
  }

  onClose(): void {
    this.dialogRef.close();
  }

  // getStatus(): ProvisionStatus {
  // let numWithAgentData = 0;
  // let totalAgents = 0;
  // for (const host of this.data.provisioningStep.ProvisionedNetworkToProvisionedHost) {
  //   totalAgents++;
  //   if (host.ProvisionedHostToAgentStatus?.clientId) numWithAgentData++;
  // }
  // if (numWithAgentData === totalAgents) {
  //   this.failedChildren = false;
  //   return ProvisionStatus.COMPLETE;
  // } else if (numWithAgentData === 0) {
  //   return ProvisionStatus.FAILED;
  // } else {
  //   this.failedChildren = true;
  //   return ProvisionStatus.INPROGRESS;
  // }
  // }

  getStatusText(): string {
    if (!this.data.planStatus) return 'Unknown';
    switch (this.data.planStatus.state) {
      case LaForgeProvisionStatus.Complete:
        return 'Complete';
      case LaForgeProvisionStatus.Failed:
        return 'Failed';
      case LaForgeProvisionStatus.Inprogress:
        return 'In Progress';
      case LaForgeProvisionStatus.Tainted:
        return 'Tainted';
      default:
        return 'Unknown';
    }
  }

  getStatusIcon(): string {
    if (!this.data.planStatus) return 'minus-circle';
    switch (this.data.planStatus.state) {
      case LaForgeProvisionStatus.Complete:
        return 'check-circle';
      case LaForgeProvisionStatus.Failed:
        return 'times-circle';
      case LaForgeProvisionStatus.Inprogress:
        return 'play-circle';
      case LaForgeProvisionStatus.Tainted:
        return 'skull';
      default:
        return 'minus-circle';
    }
  }

  getStatusColor(): string {
    if (!this.data.planStatus) return 'dark';
    switch (this.data.planStatus.state) {
      case LaForgeProvisionStatus.Complete:
        return 'success';
      case LaForgeProvisionStatus.Failed:
        return 'danger';
      case LaForgeProvisionStatus.Inprogress:
        return 'info';
      case LaForgeProvisionStatus.Tainted:
        return 'danger';
      default:
        return 'minus-circle';
    }
  }

  getText(): string {
    switch (this.data.provisioningStep.type) {
      case LaForgeProvisioningStepType.Script:
        return `${
          this.data.provisioningStep.ProvisioningStepToScript.source
        } ${this.data.provisioningStep.ProvisioningStepToScript.args.join(' ')}`;
      case LaForgeProvisioningStepType.Command:
        return `${
          this.data.provisioningStep.ProvisioningStepToCommand.program
        } ${this.data.provisioningStep.ProvisioningStepToCommand.args.join(' ')}`;
      case LaForgeProvisioningStepType.DnsRecord:
        return 'DNSRecord';
      case LaForgeProvisioningStepType.FileDownload:
        // eslint-disable-next-line max-len
        return `${this.data.provisioningStep.ProvisioningStepToFileDownload.source} -> ${this.data.provisioningStep.ProvisioningStepToFileDownload.destination}`;
      case LaForgeProvisioningStepType.FileDelete:
        return 'FileDelete';
      case LaForgeProvisioningStepType.FileExtract:
        return 'FileExtract';
      default:
        return 'Step';
    }
  }
}
