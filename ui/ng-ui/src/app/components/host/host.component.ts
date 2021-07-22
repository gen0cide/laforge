import { Component, Input, OnInit, ChangeDetectorRef, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {
  LaForgeGetBuildTreeQuery,
  LaForgeProvisionedHost,
  LaForgeProvisionStatus,
  LaForgeSubscribeUpdatedAgentStatusSubscription,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';

import { Subscription } from 'rxjs';
import { ApiService } from 'src/app/services/api/api.service';

import { RebuildService } from '../../services/rebuild/rebuild.service';
import { HostModalComponent } from '../host-modal/host-modal.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent implements OnInit, OnDestroy {
  // @Input() status: Status;
  @Input()
  // eslint-disable-next-line max-len
  provisionedHost: LaForgeGetBuildTreeQuery['build']['buildToTeam'][0]['TeamToProvisionedNetwork'][0]['ProvisionedNetworkToProvisionedHost'][0];
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() parentSelected: boolean;
  @Input() hasAgent: boolean;
  unsubscribe: Subscription[] = [];
  isSelectedState = false;
  planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  provisionedHostStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  agentStatus: LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'];
  expandOverride = false;

  constructor(
    public dialog: MatDialog,
    private rebuild: RebuildService,
    private api: ApiService,
    private envService: EnvironmentService,
    private cdRef: ChangeDetectorRef
  ) {
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
    if (!this.parentSelected) this.parentSelected = false;
    if (!this.hasAgent) this.hasAgent = false;
  }

  ngOnInit() {
    const sub1 = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkPlanStatus();
      this.checkProvisionedHostStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub1);
    const sub2 = this.envService.agentStatusUpdate.asObservable().subscribe(() => {
      this.checkAgentStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub2);
  }

  ngOnDestroy() {
    this.unsubscribe.forEach((s) => s.unsubscribe());
  }

  viewDetails(): void {
    this.dialog.open(HostModalComponent, {
      width: '50%',
      height: '80%',
      data: {
        provisionedHost: this.provisionedHost,
        planStatus: this.planStatus,
        provisionedHostStatus: this.provisionedHostStatus,
        agentStatus: this.agentStatus
      }
    });
  }

  checkPlanStatus(): void {
    this.planStatus = this.envService.getStatus(this.provisionedHost.ProvisionedHostToPlan.PlanToStatus.id) || this.planStatus;
  }

  checkProvisionedHostStatus(): void {
    this.provisionedHostStatus = this.envService.getStatus(this.provisionedHost.ProvisionedHostToStatus.id) || this.provisionedHostStatus;
  }

  checkAgentStatus(): void {
    // if (this.provisionedHost.ProvisionedHostToAgentStatus?.clientId) {
    const updatedStatus = this.envService.getAgentStatus(this.provisionedHost.id);
    if (updatedStatus) {
      // console.log('agent status updated');
      this.agentStatus = updatedStatus;
    }
    // }
  }

  isAgentStale(): boolean {
    if (!this.agentStatus) return true;
    return Date.now() / 1000 - this.agentStatus.timestamp > 120;
  }

  getStatusIcon(): string {
    const status = this.planStatus ?? this.provisionedHostStatus;
    if (status?.state) {
      switch (status.state) {
        case LaForgeProvisionStatus.Todelete:
          return 'recycle';
        case LaForgeProvisionStatus.Deleteinprogress:
          return 'trash-restore';
        case LaForgeProvisionStatus.Deleted:
          return 'trash';
      }
    }
    if (this.agentStatus) {
      if (this.isAgentStale()) return 'exclamation-circle';
      else return 'check-circle';
    } else {
      if (!status?.state) {
        return 'minus-circle';
      }
      switch (status.state) {
        case LaForgeProvisionStatus.Complete:
          return 'box-check';
        case LaForgeProvisionStatus.Failed:
          return 'times-circle';
        case LaForgeProvisionStatus.Inprogress:
          return 'play-circle';
        case LaForgeProvisionStatus.Awaiting:
          return 'spinner fa-spin';
        case LaForgeProvisionStatus.Planning:
          return 'ruler-triangle';
        default:
          return 'minus-circle';
      }
    }
  }

  getStatusColor(): string {
    const status = this.planStatus ?? this.provisionedHostStatus;
    if (status?.state) {
      switch (status.state) {
        case LaForgeProvisionStatus.Todelete:
          return 'primary';
        case LaForgeProvisionStatus.Deleteinprogress:
          return 'info';
        case LaForgeProvisionStatus.Deleted:
          return 'dark';
      }
    }
    if (this.agentStatus) {
      if (this.isAgentStale()) return 'warning';
      else return 'success';
    } else {
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
        case LaForgeProvisionStatus.Planning:
          return 'primary';
        default:
          return 'dark';
      }
    }
  }

  onSelect(): void {
    let success = false;
    if (!this.isSelected()) {
      success = this.rebuild.addHost(this.provisionedHost as LaForgeProvisionedHost);
    } else {
      success = this.rebuild.removeHost(this.provisionedHost as LaForgeProvisionedHost);
    }
    // console.log(success);
    if (success) this.isSelectedState = !this.isSelectedState;
  }

  onIndeterminateChange(isIndeterminate: boolean): void {
    if (!isIndeterminate && this.isSelectedState)
      setTimeout(() => this.rebuild.addHost(this.provisionedHost as LaForgeProvisionedHost), 500);
  }

  isSelected(): boolean {
    return this.rebuild.hasHost(this.provisionedHost as LaForgeProvisionedHost);
  }

  shouldCollapse(): boolean {
    if (this.planStatus && this.planStatus.state === LaForgeProvisionStatus.Deleted) return true;
    // return hostChildrenCompleted(this.provisionedHost as LaForgeProvisionedHost, this.envService.getStatus);
    let numCompleted = 0;
    let numAwaiting = 0;
    let totalSteps = 0;
    for (const step of this.provisionedHost.ProvisionedHostToProvisioningStep) {
      if (step.step_number === 0) continue;
      totalSteps++;
      const stepStatus = this.envService.getStatus(step.ProvisioningStepToPlan.PlanToStatus.id);
      if (stepStatus?.state === LaForgeProvisionStatus.Complete) numCompleted++;
      if (stepStatus?.state === LaForgeProvisionStatus.Awaiting) numAwaiting++;
      // if (step.ProvisioningStepToStatus.state === LaForgeProvisionStatus.Complete) numCompleted++;
    }
    // console.log(this.provisionedHost.ProvisionedHostToHost.hostname, totalSteps, numCompleted);
    if (numCompleted === totalSteps) return true;
    if (numAwaiting === totalSteps) return true;
    else return false;
  }

  toggleCollapse(): void {
    this.expandOverride = !this.expandOverride;
  }
}
