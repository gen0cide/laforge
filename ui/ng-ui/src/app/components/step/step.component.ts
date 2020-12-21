import { Component, Input } from '@angular/core';
import { ProvisionStatus } from 'src/app/models/common.model';
import { ProvisionedStep } from 'src/app/models/host.model';

@Component({
  selector: 'app-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class StepComponent {
  @Input() stepNumber: number;
  @Input() provisionedStep: ProvisionedStep;
  @Input() showDetail: boolean;
  @Input() style: 'compact' | 'expanded';

  getStatusIcon(): string {
    switch (this.provisionedStep.provisionType) {
      case 'Script':
        return 'file-code';
      case 'Command':
        return 'terminal';
      case 'DNSRecord':
        return 'globe';
      case 'FileDownload':
        return 'download';
      case 'FileDelete':
        return 'trash';
      case 'FileExtract':
        return 'file-archive';
      default:
        return 'minus-circle';
    }
  }

  getStatusColor(): string {
    switch (this.provisionedStep.status.state) {
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
