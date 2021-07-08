import { Component, Inject, OnDestroy } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ProvisionStatus } from 'src/app/models/common.model';
import { ProvisionedNetwork } from 'src/app/models/network.model';

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
    @Inject(MAT_DIALOG_DATA) public data: { provisionedNetwork: ProvisionedNetwork }
  ) {}

  onClose(): void {
    this.dialogRef.close();
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
    return ProvisionStatus[this.getStatus()];
  }

  getStatusIcon(): string {
    this.getStatus();
    if (this.failedChildren) {
      return 'skull-crossbones';
    }
    switch (this.getStatus()) {
      case ProvisionStatus.COMPLETE:
        return 'check-circle';
      case ProvisionStatus.FAILED:
        return 'times-circle';
      case ProvisionStatus.INPROGRESS:
        return 'play-circle';
      default:
        return 'minus-circle';
    }
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
}
