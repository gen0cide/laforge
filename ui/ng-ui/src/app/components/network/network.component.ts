import { ChangeDetectorRef, Component, ElementRef, Input, OnInit, Renderer2, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { ProvisionedHost } from 'src/app/models/host.model';
import { ProvisionedNetwork } from 'src/app/models/network.model';
import { NetworkModalComponent } from '../network-modal/network-modal.component';

@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss']
})
export class NetworkComponent implements OnInit {
  @Input() provisionedNetwork: ProvisionedNetwork;
  @Input() status: Status;
  // @ViewChild('options') options: ElementRef;
  optionsToggled: boolean;

  constructor(public dialog: MatDialog) {
    // private changeDetectorRef: ChangeDetectorRef // private renderer: Renderer2,
    this.optionsToggled = false;
  }

  ngOnInit(): void {
    // this.renderer.listen('window', 'click', (e: Event) => {
    //   if (!this.options.nativeElement.contains(e.target)) {
    //     this.optionsToggled = false;
    //     this.changeDetectorRef.markForCheck();
    //   }
    // });
  }

  viewDetails(): void {
    this.dialog.open(NetworkModalComponent, {
      width: '50%',
      height: '80%',
      data: { provisionedNetwork: this.provisionedNetwork }
    });
  }

  toggleOptions(): void {
    this.optionsToggled = !this.optionsToggled;
  }

  getStatus(): ProvisionStatus {
    let status: ProvisionStatus = this.status.state;
    for (const host of this.provisionedNetwork.provisionedHosts) {
      switch (host.status.state) {
        case ProvisionStatus.ProvStatusFailed:
          status = ProvisionStatus.ProvStatusFailed;
          break;
        case ProvisionStatus.ProvStatusInProgress:
          if (status === ProvisionStatus.ProvStatusComplete) status = ProvisionStatus.ProvStatusInProgress;
          break;
        default:
          break;
      }
    }
    return status;
  }

  getStatusColor(): string {
    switch (this.getStatus()) {
      case ProvisionStatus.ProvStatusComplete:
        return 'success';
      case ProvisionStatus.ProvStatusInProgress:
        return 'info';
      case ProvisionStatus.ProvStatusFailed:
        return 'danger';
      default:
        return 'dark';
    }
  }
}
