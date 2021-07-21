import { ChangeDetectorRef, Component, Inject, OnInit, OnDestroy } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import {
  LaForgeSubscribeUpdatedStatusSubscription,
  LaForgeSubscribeUpdatedAgentStatusSubscription,
  LaForgeGetBuildTreeQuery,
  LaForgeProvisionStatus
} from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';
import { ApiService } from 'src/app/services/api/api.service';

@Component({
  selector: 'app-host-modal',
  templateUrl: './host-modal.component.html',
  styleUrls: ['./host-modal.component.scss']
})
class HostModalComponent implements OnInit, OnDestroy {
  varsColumns: string[] = ['key', 'value'];
  tagsColumns: string[] = ['name', 'description'];
  // planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  // provisionedHostStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  // agentStatus: LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'];

  constructor(
    public dialogRef: MatDialogRef<HostModalComponent>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      // eslint-disable-next-line max-len
      provisionedHost: LaForgeGetBuildTreeQuery['build']['buildToTeam'][0]['TeamToProvisionedNetwork'][0]['ProvisionedNetworkToProvisionedHost'][0];
      planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
      provisionedHostStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
      agentStatus: LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'];
      needsToQuerySteps?: boolean;
    },
    private api: ApiService,
    private envService: EnvironmentService,
    private cdRef: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    // this.envService.statusUpdate.asObservable().subscribe(() => {
    //   this.checkPlanStatus();
    //   this.checkProvisionedHostStatus();
    //   this.cdRef.detectChanges();
    // });
    // this.envService.agentStatusUpdate.asObservable().subscribe(() => {
    //   this.checkAgentStatus();
    //   this.cdRef.detectChanges();
    // });
  }

  ngOnDestroy() {
    // this.envService.statusUpdate.unsubscribe();
    // this.envService.agentStatusUpdate.unsubscribe();
  }

  onClose(): void {
    this.dialogRef.close();
  }

  checkPlanStatus(): void {
    this.data.planStatus =
      this.envService.getStatus(this.data.provisionedHost.ProvisionedHostToPlan.PlanToStatus.id) || this.data.planStatus;
  }

  checkProvisionedHostStatus(): void {
    this.data.provisionedHostStatus =
      this.envService.getStatus(this.data.provisionedHost.ProvisionedHostToStatus.id) || this.data.provisionedHostStatus;
  }

  checkAgentStatus(): void {
    this.data.agentStatus =
      this.envService.getAgentStatus(this.data.provisionedHost.ProvisionedHostToAgentStatus.clientId) || this.data.agentStatus;
  }

  isAgentStale(): boolean {
    if (!this.data.agentStatus) return true;
    return Date.now() / 1000 - this.data.agentStatus.timestamp > 120;
  }

  getStatusIcon(): string {
    if (this.data.agentStatus) {
      if (this.isAgentStale()) return 'exclamation-circle';
      else return 'check-circle';
    }
    const status = this.data.planStatus ?? this.data.provisionedHostStatus;
    if (!status?.state) {
      return 'minus-circle';
    }
    switch (status.state) {
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
    if (this.data.agentStatus) {
      if (this.isAgentStale()) return 'warning';
      else return 'success';
    }
    const status = this.data.planStatus ?? this.data.provisionedHostStatus;
    if (!status?.state) {
      return 'minus-circle';
    }
    switch (status.state) {
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

export { HostModalComponent };
