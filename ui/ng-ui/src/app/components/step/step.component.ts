import { Component, Input } from '@angular/core';
import { ProvisionStatus } from 'src/app/models/common.model';
import { ProvisioningStep, ProvisioningStepType } from 'src/app/models/step.model';

@Component({
  selector: 'app-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class StepComponent {
  @Input() stepNumber: number;
  @Input() provisionedStep: ProvisioningStep;
  @Input() showDetail: boolean;
  @Input() style: 'compact' | 'expanded';

  getStatusIcon(): string {
    switch (this.provisionedStep.type) {
      case ProvisioningStepType.Script:
        return 'file-code';
      case ProvisioningStepType.Command:
        return 'terminal';
      case ProvisioningStepType.DNSRecord:
        return 'globe';
      case ProvisioningStepType.FileDownload:
        return 'download';
      case ProvisioningStepType.FileDelete:
        return 'trash';
      case ProvisioningStepType.FileExtract:
        return 'file-archive';
      default:
        return 'minus-circle';
    }
  }

  getStatusColor(): string {
    // const status = this.provisionedStep.ProvisioningStepToPlan?.PlanToStatus ?? this.provisionedStep.ProvisioningStepToStatus;
    const status = this.provisionedStep.ProvisioningStepToStatus;
    switch (status.state) {
      case ProvisionStatus.COMPLETE:
        return 'success';
      case ProvisionStatus.TODELETE:
        return 'warning';
      case ProvisionStatus.DELETEINPROGRESS:
        return 'warning';
      case ProvisionStatus.DELETED:
        return 'dark';
      case ProvisionStatus.FAILED:
        return 'danger';
      case ProvisionStatus.INPROGRESS:
        return 'info';
      default:
        return 'dark';
    }
  }

  getText(): string {
    switch (this.provisionedStep.type) {
      case ProvisioningStepType.Script:
        return `${this.provisionedStep.ProvisioningStepToScript.source} ${this.provisionedStep.ProvisioningStepToScript.args.join(' ')}`;
      case ProvisioningStepType.Command:
        return `${this.provisionedStep.ProvisioningStepToCommand.program} ${this.provisionedStep.ProvisioningStepToCommand.args.join(' ')}`;
      case ProvisioningStepType.DNSRecord:
        return 'DNSRecord';
      case ProvisioningStepType.FileDownload:
        // eslint-disable-next-line max-len
        return `${this.provisionedStep.ProvisioningStepToFileDownload.source} -> ${this.provisionedStep.ProvisioningStepToFileDownload.destination}`;
      case ProvisioningStepType.FileDelete:
        return 'FileDelete';
      case ProvisioningStepType.FileExtract:
        return 'FileExtract';
      default:
        return 'Step';
    }
  }
}
