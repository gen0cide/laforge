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
  @Input() style: 'compact' | 'collapsed' | 'expanded';

  constructor(public dialog: MatDialog) {
    if (!this.style) this.style = 'compact';
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

  // toggleOptions(): void {
  //   this.optionsToggled = !this.optionsToggled;
  // }

  getStatus(): ProvisionStatus {
    let numWithAgentData = 0;
    let totalAgents = 0;
    for (const host of this.provisionedNetwork.provisionedHosts) {
      totalAgents++;
      if (host.heartbeat) numWithAgentData++;
    }
    if (numWithAgentData === totalAgents) return ProvisionStatus.ProvStatusComplete;
    else if (numWithAgentData === 0) return ProvisionStatus.ProvStatusFailed;
    else return ProvisionStatus.ProvStatusInProgress;
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
