import { Component, Input } from '@angular/core';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { MatDialog } from '@angular/material/dialog';
import { ProvisionedHost } from 'src/app/models/host.model';
import { HostModalComponent } from '../host-modal/host-modal.component';
import { RebuildService } from '../../services/rebuild/rebuild.service';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent {
  @Input() status: Status;
  @Input() provisionedHost: ProvisionedHost;
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() parentSelected: boolean;
  isSelected = false;

  constructor(public dialog: MatDialog, private rebuild: RebuildService) {
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
    if (!this.parentSelected) this.parentSelected = false;
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

  onSelect(): void {
    this.isSelected = !this.isSelected;
    if (this.isSelected) this.rebuild.addHost(this.provisionedHost);
    else this.rebuild.removeHost(this.provisionedHost);
  }

  onIndeterminateChange(isIndeterminate: boolean): void {
    if (!isIndeterminate && this.isSelected) setTimeout(() => this.rebuild.addHost(this.provisionedHost), 500);
  }
}
