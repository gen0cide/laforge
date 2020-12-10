import { ChangeDetectorRef, Component, ElementRef, Input, OnInit, Renderer2, ViewChild } from '@angular/core';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { ProvisionedHost } from 'src/app/models/host.model';

@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss']
})
export class NetworkComponent implements OnInit {
  @Input() hosts: ProvisionedHost[];
  @Input() title: string;
  @Input() details: string;
  @Input() status: Status;
  // @ViewChild('options') options: ElementRef;
  optionsToggled: boolean;

  constructor() {
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

  toggleOptions(): void {
    this.optionsToggled = !this.optionsToggled;
  }

  getStatus(): ProvisionStatus {
    let status: ProvisionStatus = this.status.state;
    for (const host of this.hosts) {
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
