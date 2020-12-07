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
import { ProvisionedNetwork } from 'src/app/models/network.model';

@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss']
})
export class NetworkComponent implements OnInit {
  @Input() provisionedNetwork: ProvisionedNetwork;
  @ViewChild('options') options: ElementRef;
  optionsToggled: boolean;

  constructor(
    private renderer: Renderer2,
    private changeDetectorRef: ChangeDetectorRef
  ) {
    this.optionsToggled = false;
  }

  ngOnInit(): void {
    // console.log(this.provisionedNetwork);
    this.renderer.listen('window', 'click', (e: Event) => {
      if (!this.options.nativeElement.contains(e.target)) {
        this.optionsToggled = false;
        this.changeDetectorRef.markForCheck();
      }
    });
  }

  toggleOptions(): void {
    this.optionsToggled = !this.optionsToggled;
  }

  getStatus(): string {
    switch (this.provisionedNetwork.status.state) {
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
