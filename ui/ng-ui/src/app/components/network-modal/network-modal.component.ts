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
  tagsColumns: string[] = ['name', 'description'];
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
    for (const host of this.data.provisionedNetwork.provisionedHosts) {
      totalAgents++;
      if (host.heartbeat) numWithAgentData++;
    }
    if (numWithAgentData === totalAgents) return ProvisionStatus.ProvStatusComplete;
    else if (numWithAgentData === 0) return ProvisionStatus.ProvStatusFailed;
    else return ProvisionStatus.ProvStatusInProgress;
  }

  getStatusText(): string {
    return ProvisionStatus[this.getStatus()];
  }

  getStatusIcon(): string {
    switch (this.getStatus()) {
      case ProvisionStatus.ProvStatusComplete:
        return 'check-circle';
      case ProvisionStatus.ProvStatusFailed:
        return 'times-circle';
      case ProvisionStatus.ProvStatusInProgress:
        return 'play-circle';
      default:
        return 'minus-circle';
    }
  }

  getStatusColor(): string {
    switch (this.getStatus()) {
      case ProvisionStatus.ProvStatusComplete:
        return 'success';
      case ProvisionStatus.ProvStatusInProgress:
        return 'warning';
      case ProvisionStatus.ProvStatusFailed:
        return 'danger';
      default:
        return 'dark';
    }
  }
}
