import { ChangeDetectorRef, Component, ElementRef, Input, OnInit, Renderer2, ViewChild } from '@angular/core';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { MatDialog } from '@angular/material/dialog';
import { ProvisionedHost } from 'src/app/models/host.model';
import { HostModalComponent } from '../host-modal/host-modal.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent implements OnInit {
  @Input() status: Status;
  @Input() provisionedHost: ProvisionedHost;
  // @ViewChild('container') container: ElementRef;

  constructor(public dialog: MatDialog) {
    // private changeDetectorRef: ChangeDetectorRef // private renderer: Renderer2,
  }

  ngOnInit(): void {
    // this.renderer.listen('window', 'click', (e: Event) => {
    //   if (!this.container.nativeElement.contains(e.target)) {
    //     this.optionsToggled = false;
    //     this.changeDetectorRef.markForCheck();
    //   }
    // });
  }

  viewDetails(): void {
    this.dialog.open(HostModalComponent, {
      width: '50%',
      data: { provisionedHost: this.provisionedHost }
    });
  }

  getStatusIcon(): string {
    switch (this.status.state) {
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
    switch (this.status.state) {
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
