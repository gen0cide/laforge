import { Component, Input } from '@angular/core';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { MatDialog } from '@angular/material/dialog';
import { ProvisionedHost } from 'src/app/models/host.model';
import { HostModalComponent } from '../host-modal/host-modal.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent {
  @Input() status: Status;
  @Input() provisionedHost: ProvisionedHost;
  @Input() style: 'compact' | 'collapsed' | 'expanded';

  constructor(public dialog: MatDialog) {
    if (!this.style) this.style = 'compact';
  }

  viewDetails(): void {
    this.dialog.open(HostModalComponent, {
      width: '50%',
      height: '80%',
      data: { provisionedHost: this.provisionedHost }
    });
  }

  getStatusIcon(): string {
    if (this.provisionedHost.heartbeat) return 'check';
    else return 'minus';
    // switch (this.status.state) {
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
    if (this.provisionedHost.heartbeat) return 'success';
    else return 'dark';
    // switch (this.status.state) {
    //   case ProvisionStatus.ProvStatusComplete:
    //     return 'success';
    //   case ProvisionStatus.ProvStatusFailed:
    //     return 'danger';
    //   case ProvisionStatus.ProvStatusInProgress:
    //     return 'info';
    //   default:
    //     return 'dark';
    // }
  }
}
