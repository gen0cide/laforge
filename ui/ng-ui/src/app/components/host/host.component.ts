import {
  ChangeDetectorRef,
  Component,
  ElementRef,
  Input,
  OnInit,
  Renderer2,
  ViewChild
} from '@angular/core';
import { ProvisionStatus } from 'src/app/models/common.model';
import { ProvisionedHost } from 'src/app/models/host.model';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent implements OnInit {
  @Input() provisionedHost: ProvisionedHost;
  @ViewChild('container') container: ElementRef;
  optionsToggled: boolean;

  constructor(
    private renderer: Renderer2,
    private changeDetectorRef: ChangeDetectorRef
  ) {
    this.optionsToggled = false;
  }

  ngOnInit(): void {
    this.renderer.listen('window', 'click', (e: Event) => {
      if (!this.container.nativeElement.contains(e.target)) {
        this.optionsToggled = false;
        this.changeDetectorRef.markForCheck();
      }
    });
  }

  toggleOptions(): void {
    this.optionsToggled = !this.optionsToggled;
  }

  getStatus(): string {
    switch (this.provisionedHost.status.state) {
      case ProvisionStatus.ProvStatusComplete:
        return 'ProvStatusComplete';
      case ProvisionStatus.ProvStatusFailed:
        return 'ProvStatusFailed';
      case ProvisionStatus.ProvStatusInProgress:
        return 'ProvStatusInProgress';
      default:
        return 'ProvStatusUndefined';
    }
  }
}
