import { Component, Inject } from '@angular/core';
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

  getStatusIcon(): string {
    switch (this.data.provisionedNetwork.status.state) {
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
    switch (this.data.provisionedNetwork.status.state) {
      case ProvisionStatus.ProvStatusComplete:
        return 'success';
      case ProvisionStatus.ProvStatusFailed:
        return 'danger';
      case ProvisionStatus.ProvStatusInProgress:
        return 'info';
      default:
        return 'dark';
    }
  }
}
