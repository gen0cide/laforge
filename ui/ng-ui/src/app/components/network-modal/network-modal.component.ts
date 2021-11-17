import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { LaForgeGetBuildTreeQuery, LaForgeProvisionStatus, LaForgeSubscribeUpdatedStatusSubscription } from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';
import { ProvisionStatus } from 'src/app/models/common.model';

@Component({
  selector: 'app-network-modal',
  templateUrl: './network-modal.component.html',
  styleUrls: ['./network-modal.component.scss']
})
export class NetworkModalComponent {
  varsColumns: string[] = ['key', 'value'];
  tagsColumns: string[] = ['key', 'value'];
  failedChildren = false;

  constructor(
    public dialogRef: MatDialogRef<NetworkModalComponent>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      provisionedNetwork: LaForgeGetBuildTreeQuery['build']['buildToTeam'][0]['TeamToProvisionedNetwork'][0];
      planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
    },
    private envService: EnvironmentService
  ) {}

  onClose(): void {
    this.dialogRef.close();
  }

  allChildrenResponding(): boolean {
    let numWithAgentData = 0;
    let numWithCompletedSteps = 0;
    let totalHosts = 0;
    for (const host of this.data.provisionedNetwork.ProvisionedNetworkToProvisionedHost) {
      totalHosts++;
      if (host.ProvisionedHostToAgentStatus?.clientId) numWithAgentData++;
      let totalSteps = 0;
      let totalCompletedSteps = 0;
      for (const step of host.ProvisionedHostToProvisioningStep) {
        if (step.step_number === 0) continue;
        totalSteps++;
        if (
          step.ProvisioningStepToStatus.id &&
          this.envService.getStatus(step.ProvisioningStepToPlan.PlanToStatus.id)?.state === LaForgeProvisionStatus.Complete
        )
          totalCompletedSteps++;
      }
      if (totalSteps === totalCompletedSteps) numWithCompletedSteps++;
    }
    return numWithAgentData === totalHosts && numWithCompletedSteps === totalHosts;
  }

  getStatus(): ProvisionStatus {
    let numWithAgentData = 0;
    let totalAgents = 0;
    for (const host of this.data.provisionedNetwork.ProvisionedNetworkToProvisionedHost) {
      totalAgents++;
      if (host.ProvisionedHostToAgentStatus?.clientId) numWithAgentData++;
    }
    if (numWithAgentData === totalAgents) {
      this.failedChildren = false;
      return ProvisionStatus.COMPLETE;
    } else if (numWithAgentData === 0) {
      return ProvisionStatus.FAILED;
    } else {
      this.failedChildren = true;
      return ProvisionStatus.INPROGRESS;
    }
  }

  getStatusText(): string {
    switch (this.data.planStatus.state) {
      case LaForgeProvisionStatus.Complete:
        return 'Complete';
      case LaForgeProvisionStatus.Failed:
        return 'Failed';
      case LaForgeProvisionStatus.Inprogress:
        return 'In Progress';
      default:
        return 'minus-circle';
    }
  }

  getStatusIcon(): string {
    if (this.failedChildren) {
      return 'skull-crossbones';
    }
    switch (this.data.planStatus.state) {
      case LaForgeProvisionStatus.Complete:
        return 'check-circle';
      case LaForgeProvisionStatus.Failed:
        return 'times-circle';
      case LaForgeProvisionStatus.Inprogress:
        return 'play-circle';
      default:
        return 'minus-circle';
    }
  }

  getStatusColor(): string {
    if (this.failedChildren) {
      return 'warning';
    }
    switch (this.data.planStatus.state) {
      case LaForgeProvisionStatus.Complete:
        return 'success';
      case LaForgeProvisionStatus.Failed:
        return 'danger';
      case LaForgeProvisionStatus.Inprogress:
        return 'info';
      default:
        return 'dark';
    }
  }
}
