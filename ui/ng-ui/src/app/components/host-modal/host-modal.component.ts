import { Component, Inject, Input, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ProvisionStatus } from 'src/app/models/common.model';
import { ProvisionedHost } from 'src/app/models/host.model';
import { ProvisioningStep } from 'src/app/models/step.model';
import { ApiService } from 'src/app/services/api/api.service';

@Component({
  selector: 'app-host-modal',
  templateUrl: './host-modal.component.html',
  styleUrls: ['./host-modal.component.scss']
})
class HostModalComponent implements OnInit {
  varsColumns: string[] = ['key', 'value'];
  tagsColumns: string[] = ['name', 'description'];
  provisionedSteps: ProvisioningStep[];

  constructor(
    public dialogRef: MatDialogRef<HostModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { provisionedHost: ProvisionedHost; needsToQuerySteps?: boolean },
    private api: ApiService
  ) {}

  ngOnInit(): void {
    if (this.data.provisionedHost.ProvisionedHostToProvisioningStep.length === 0) {
      this.api.pullHostSteps(this.data.provisionedHost.id).then((steps: ProvisioningStep[]) => {
        this.provisionedSteps = steps;
      });
    }
  }

  onClose(): void {
    this.dialogRef.close();
  }

  isAgentStale(): boolean {
    if (!this.data.provisionedHost.ProvisionedHostToAgentStatus?.clientId) return true;
    return Date.now() / 1000 - this.data.provisionedHost.ProvisionedHostToAgentStatus.timestamp > 60;
  }

  getStatusIcon(): string {
    if (this.data.provisionedHost.ProvisionedHostToAgentStatus?.clientId) {
      if (this.isAgentStale()) return 'exclamation-circle';
      else return 'check-circle';
    }
    // TODO: Fix for live statuses after finals
    switch (this.data.provisionedHost.ProvisionedHostToStatus.state) {
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
    if (this.data.provisionedHost.ProvisionedHostToAgentStatus?.clientId) {
      if (this.isAgentStale()) return 'warning';
      else return 'success';
    }
    // TODO: Fix for live statuses after finals
    switch (this.data.provisionedHost.ProvisionedHostToStatus.state) {
      case ProvisionStatus.COMPLETE:
        return 'success';
      case ProvisionStatus.FAILED:
        return 'danger';
      case ProvisionStatus.INPROGRESS:
        return 'info';
      default:
        return 'dark';
    }
  }
}

export { HostModalComponent };
