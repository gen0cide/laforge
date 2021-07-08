import { Component, Input } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { hostChildrenCompleted, ProvisionedHost } from 'src/app/models/host.model';

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
  @Input() provisionedHost: ProvisionedHost;
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() parentSelected: boolean;
  @Input() hasAgent: boolean;
  isSelectedState = false;

  constructor(public dialog: MatDialog, private rebuild: RebuildService, private api: ApiService) {
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
    if (!this.parentSelected) this.parentSelected = false;
    if (!this.hasAgent) this.hasAgent = false;
  }

  viewDetails(): void {
    this.dialog.open(HostModalComponent, {
      width: '50%',
      height: '80%',
      data: { provisionedHost: this.provisionedHost }
    });
  }

  isAgentStale(): boolean {
    if (!this.provisionedHost.ProvisionedHostToAgentStatus?.clientId) return true;
    return Date.now() / 1000 - this.provisionedHost.ProvisionedHostToAgentStatus.timestamp > 60;
  }

  getStatusIcon(): string {
    if (this.provisionedHost.ProvisionedHostToAgentStatus?.clientId) {
      if (this.isAgentStale()) return 'exclamation-circle';
      else return 'check-circle';
    } else {
      const status = this.provisionedHost.ProvisionedHostToPlan?.PlanToStatus ?? this.provisionedHost.ProvisionedHostToStatus;
      switch (status.state) {
        case ProvisionStatus.COMPLETE:
          return 'box-check';
        case ProvisionStatus.TODELETE:
          return 'recycle';
        case ProvisionStatus.DELETEINPROGRESS:
          return 'trash-restore';
        case ProvisionStatus.DELETED:
          return 'trash';
        case ProvisionStatus.FAILED:
          return 'times-circle';
        case ProvisionStatus.INPROGRESS:
          return 'play-circle';
        default:
          return 'minus-circle';
      }
    }
  }

  getStatusColor(): string {
    if (this.provisionedHost.ProvisionedHostToAgentStatus?.clientId) {
      if (this.isAgentStale()) return 'warning';
      else return 'success';
    } else {
      const status = this.provisionedHost.ProvisionedHostToPlan?.PlanToStatus ?? this.provisionedHost.ProvisionedHostToStatus;
      switch (status.state) {
        case ProvisionStatus.COMPLETE:
          return 'success';
        case ProvisionStatus.TODELETE:
          return 'recycle';
        case ProvisionStatus.DELETEINPROGRESS:
          return 'trash-restore';
        case ProvisionStatus.DELETED:
          return 'trash';
        case ProvisionStatus.FAILED:
          return 'danger';
        case ProvisionStatus.INPROGRESS:
          return 'info';
        default:
          return 'dark';
      }
    }
  }

  onSelect(): void {
    let success = false;
    if (!this.isSelected()) {
      success = this.rebuild.addHost(this.provisionedHost);
    } else {
      success = this.rebuild.removeHost(this.provisionedHost);
    }
    console.log(success);
    if (success) this.isSelectedState = !this.isSelectedState;
  }

  onIndeterminateChange(isIndeterminate: boolean): void {
    if (!isIndeterminate && this.isSelectedState) setTimeout(() => this.rebuild.addHost(this.provisionedHost), 500);
  }

  isSelected(): boolean {
    return this.rebuild.hasHost(this.provisionedHost);
  }

  shouldCollapse(): boolean {
    return hostChildrenCompleted(this.provisionedHost);
  }
}
