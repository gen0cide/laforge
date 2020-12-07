import {
  ChangeDetectorRef,
  Component,
  ElementRef,
  Input,
  OnInit,
  Renderer2,
  ViewChild
} from '@angular/core';
import { ProvisionedHost } from 'src/app/models/host.model';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent implements OnInit {
  @Input() provisionedHost: ProvisionedHost;
  @ViewChild('options') options: ElementRef;
  optionsToggled: boolean;

  constructor(
    private renderer: Renderer2,
    private changeDetectorRef: ChangeDetectorRef
  ) {
    this.optionsToggled = false;
  }

  ngOnInit(): void {
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
}
