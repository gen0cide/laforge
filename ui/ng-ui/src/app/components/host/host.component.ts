import { Component, Input } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {
  LaForgeGetBuildTreeQuery,
  LaForgeProvisionedHost,
  LaForgeProvisionStatus,
  LaForgeSubscribeUpdatedAgentStatusSubscription,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';
import { hostChildrenCompleted } from '@util';

import { ApiService } from 'src/app/services/api/api.service';

import { RebuildService } from '../../services/rebuild/rebuild.service';
import { HostModalComponent } from '../host-modal/host-modal.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent {
  // @Input() status: Status;
  @Input()
  // eslint-disable-next-line max-len
  provisionedHost: LaForgeGetBuildTreeQuery['build']['buildToTeam'][0]['TeamToProvisionedNetwork'][0]['ProvisionedNetworkToProvisionedHost'][0];
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() parentSelected: boolean;
  @Input() hasAgent: boolean;
  isSelectedState = false;
  // planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  // provisionedHostStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  // agentStatus: LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'];

  constructor(public dialog: MatDialog, private rebuild: RebuildService, private api: ApiService, private envService: EnvironmentService) {
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
    if (!this.parentSelected) this.parentSelected = false;
    if (!this.hasAgent) this.hasAgent = false;

    // envService.statusUpdate.subscribe(() => {

    // })
  }

  viewDetails(): void {
    this.dialog.open(HostModalComponent, {
      width: '50%',
      height: '80%',
      data: { provisionedHost: this.provisionedHost }
    });
  }

  getPlanStatus(): LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'] {
    return this.envService.getStatus(this.provisionedHost.ProvisionedHostToPlan.id);
  }

  getProvisionedHostStatus(): LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'] {
    return this.envService.getStatus(this.provisionedHost.ProvisionedHostToStatus.id);
  }

  getAgentStatus(): LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'] {
    return this.envService.getAgentStatus(this.provisionedHost.ProvisionedHostToAgentStatus.clientId);
  }

  isAgentStale(): boolean {
    if (!this.getAgentStatus()) return true;
    return Date.now() / 1000 - this.getAgentStatus().timestamp > 60;
  }

  getStatusIcon(): string {
    if (this.getAgentStatus()) {
      if (this.isAgentStale()) return 'exclamation-circle';
      else return 'check-circle';
    } else {
      const status = this.getPlanStatus() ?? this.getProvisionedHostStatus();
      if (!status?.state) {
        return 'minus-circle';
      }
      switch (status.state) {
        case LaForgeProvisionStatus.Complete:
          return 'box-check';
        case LaForgeProvisionStatus.Todelete:
          return 'recycle';
        case LaForgeProvisionStatus.Deleteinprogress:
          return 'trash-restore';
        case LaForgeProvisionStatus.Deleted:
          return 'trash';
        case LaForgeProvisionStatus.Failed:
          return 'times-circle';
        case LaForgeProvisionStatus.Inprogress:
          return 'play-circle';
        default:
          return 'minus-circle';
      }
    }
  }

  getStatusColor(): string {
    if (this.getAgentStatus()) {
      if (this.isAgentStale()) return 'warning';
      else return 'success';
    } else {
      const status = this.getPlanStatus() ?? this.getProvisionedHostStatus();
      if (!status?.state) {
        return 'minus-circle';
      }
      switch (status.state) {
        case LaForgeProvisionStatus.Complete:
          return 'success';
        case LaForgeProvisionStatus.Todelete:
          return 'recycle';
        case LaForgeProvisionStatus.Deleteinprogress:
          return 'trash-restore';
        case LaForgeProvisionStatus.Deleted:
          return 'trash';
        case LaForgeProvisionStatus.Failed:
          return 'danger';
        case LaForgeProvisionStatus.Inprogress:
          return 'info';
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
    console.log(success);
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
    return hostChildrenCompleted(this.provisionedHost as LaForgeProvisionedHost);
  }
}
