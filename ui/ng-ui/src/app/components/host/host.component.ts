import {
  ChangeDetectorRef,
  Component,
  ElementRef,
  Input,
  OnInit,
  Renderer2,
  ViewChild
} from '@angular/core';
import { ProvisionStatus, Status } from 'src/app/models/common.model';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent implements OnInit {
  @Input() status: Status;
  @Input() title: string;
  // @ViewChild('container') container: ElementRef;

  constructor() {
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
