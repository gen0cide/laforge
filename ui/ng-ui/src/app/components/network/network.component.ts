import { Component, Input, OnInit } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialog } from '@angular/material/dialog';
import { LaForgeProvisionedNetwork } from '@graphql';
import { networkChildrenCompleted } from '@util';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { RebuildService } from 'src/app/services/rebuild/rebuild.service';

import { NetworkModalComponent } from '../network-modal/network-modal.component';

@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss']
})
export class NetworkComponent implements OnInit {
  @Input() provisionedNetwork: LaForgeProvisionedNetwork;
  @Input() status: Status;
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() parentSelected: boolean;
  isSelectedState = false;

  constructor(public dialog: MatDialog, private rebuild: RebuildService) {
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
    if (!this.parentSelected) this.parentSelected = false;
  }

  ngOnInit(): void {}

  viewDetails(): void {
    this.dialog.open(NetworkModalComponent, {
      width: '50%',
      height: '80%',
      data: { provisionedNetwork: this.provisionedNetwork }
    });
  }

  getStatus(): ProvisionStatus {
    let numWithAgentData = 0;
    let totalAgents = 0;
    for (const host of this.provisionedNetwork.ProvisionedNetworkToProvisionedHost) {
      totalAgents++;
      if (host.ProvisionedHostToAgentStatus?.clientId) numWithAgentData++;
    }
    if (numWithAgentData === totalAgents) return ProvisionStatus.COMPLETE;
    else if (numWithAgentData === 0) return ProvisionStatus.FAILED;
    else return ProvisionStatus.INPROGRESS;
  }

  getStatusColor(): string {
    switch (this.getStatus()) {
      case ProvisionStatus.COMPLETE:
        return 'success';
      case ProvisionStatus.INPROGRESS:
        return 'warning';
      case ProvisionStatus.FAILED:
        return 'danger';
      default:
        return 'dark';
    }
  }

  onSelect(): void {
    let success = false;
    if (!this.isSelected()) {
      success = this.rebuild.addNetwork(this.provisionedNetwork);
    } else {
      success = this.rebuild.removeNetwork(this.provisionedNetwork);
    }
    if (success) this.isSelectedState = !this.isSelectedState;
  }

  onIndeterminateChange(isIndeterminate: boolean): void {
    if (!isIndeterminate && this.isSelectedState) setTimeout(() => this.rebuild.addNetwork(this.provisionedNetwork), 500);
  }

  isSelected(): boolean {
    return this.rebuild.hasNetwork(this.provisionedNetwork);
  }

  shouldCollapse(): boolean {
    return networkChildrenCompleted(this.provisionedNetwork);
  }
}
