import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ProvisionStatus } from 'src/app/models/common.model';
import { ProvisionedHost } from 'src/app/models/host.model';

@Component({
  selector: 'app-host-modal',
  templateUrl: './host-modal.component.html',
  styleUrls: ['./host-modal.component.scss']
})
class HostModalComponent {
  varsColumns: string[] = ['key', 'value'];
  tagsColumns: string[] = ['name', 'description'];
  constructor(public dialogRef: MatDialogRef<HostModalComponent>, @Inject(MAT_DIALOG_DATA) public data: { provisionedHost: ProvisionedHost }) {
    console.log(data.provisionedHost);
  }

  onClose(): void {
    this.dialogRef.close();
  }

  getStatusIcon(): string {
    if (this.data.provisionedHost.heartbeat) return 'check-circle';
    else return 'minus-circle';
    // switch (this.data.provisionedHost.status.state) {
    //   case ProvisionStatus.ProvStatusComplete:
    //     return 'check-circle';
    //   case ProvisionStatus.ProvStatusFailed:
    //     return 'times-circle';
    //   case ProvisionStatus.ProvStatusInProgress:
    //     return 'play-circle';
    //   default:
    //     return 'minus-circle';
    // }
  }

  getStatusColor(): string {
    if (this.data.provisionedHost.heartbeat) return 'success';
    else return 'dark';
    switch (this.data.provisionedHost.status.state) {
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

export { HostModalComponent };
