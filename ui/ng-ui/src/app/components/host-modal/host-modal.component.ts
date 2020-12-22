import { Component, Inject, Input, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ProvisionStatus } from 'src/app/models/common.model';
import { ProvisionedHost, ProvisionedStep } from 'src/app/models/host.model';
import { ApiService } from 'src/app/services/api/api.service';

@Component({
  selector: 'app-host-modal',
  templateUrl: './host-modal.component.html',
  styleUrls: ['./host-modal.component.scss']
})
class HostModalComponent implements OnInit {
  varsColumns: string[] = ['key', 'value'];
  tagsColumns: string[] = ['name', 'description'];
  provisionedSteps: ProvisionedStep[];

  constructor(
    public dialogRef: MatDialogRef<HostModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { provisionedHost: ProvisionedHost; needsToQuerySteps?: boolean },
    private api: ApiService
  ) {}

  ngOnInit(): void {
    if (this.data.provisionedHost.provisionedSteps == null) {
      this.api.pullHostSteps(this.data.provisionedHost.id).then((steps: ProvisionedStep[]) => {
        this.provisionedSteps = steps;
      });
    }
  }

  onClose(): void {
    this.dialogRef.close();
  }

  getStatusIcon(): string {
    if (this.data.provisionedHost.heartbeat) return 'check-circle';
    else return 'minus-circle';
    // TODO: Fix for live statuses after finals
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
    // TODO: Fix for live statuses after finals
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
