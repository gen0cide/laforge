import { Component, Input, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {
  LaForgeGetBuildTreeQuery,
  LaForgeProvisioningStepType,
  LaForgeSubscribeUpdatedStatusSubscription,
  LaForgeProvisionStatus,
  LaForgePlanFieldsFragment
} from '@graphql';
import { Subscription } from 'rxjs';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

import { StepModalComponent } from '@components/step-modal/step-modal.component';

@Component({
  selector: 'app-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class StepComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
  @Input() stepNumber: number;
  @Input()
  // eslint-disable-next-line max-len
  provisioningStep: LaForgeGetBuildTreeQuery['build']['buildToTeam'][0]['TeamToProvisionedNetwork'][0]['ProvisionedNetworkToProvisionedHost'][0]['ProvisionedHostToProvisioningStep'][0];
  @Input() showDetail: boolean;
  @Input() style: 'compact' | 'expanded';
  @Input() mode: 'plan' | 'build' | 'manage';
  planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  provisioningStepStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  latestDiff: LaForgePlanFieldsFragment['PlanToPlanDiffs'][0];

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private envService: EnvironmentService,
    private dialog: MatDialog
  ) {
    if (!this.mode) this.mode = 'manage';
  }

  ngOnInit() {
    const sub = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkPlanStatus();
      this.checkprovisioningStepStatus();
    });
    this.unsubscribe.push(sub);
    // if (this.mode === 'plan') {
    //   const sub2 = this.envService.planUpdate.asObservable().subscribe(() => {
    //     this.checkLatestPlanDiff();
    //     this.cdRef.markForCheck();
    //   });
    //   this.unsubscribe.push(sub2);
    // }
  }

  ngOnDestroy() {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
  }

  viewDetails(): void {
    if (
      this.planStatus.state === LaForgeProvisionStatus.Awaiting ||
      this.planStatus.state === LaForgeProvisionStatus.Deleted ||
      this.planStatus.state === LaForgeProvisionStatus.Planning
    )
      return;
    this.dialog.open(StepModalComponent, {
      width: '50%',
      height: '80%',
      data: {
        provisioningStep: this.provisioningStep,
        planStatus: this.planStatus
      }
    });
  }

  checkPlanStatus(): void {
    if (!this.provisioningStep.ProvisioningStepToPlan) return;
    const updatedStatus = this.envService.getStatus(this.provisioningStep.ProvisioningStepToPlan.PlanToStatus.id);
    if (updatedStatus) {
      this.planStatus = updatedStatus;
      this.cdRef.markForCheck();
    }
  }

  checkprovisioningStepStatus(): void {
    const updatedStatus = this.envService.getStatus(this.provisioningStep.ProvisioningStepToStatus.id);
    if (updatedStatus) {
      this.provisioningStepStatus = updatedStatus;
      this.cdRef.markForCheck();
    }
  }

  checkLatestPlanDiff(): void {
    if (!this.provisioningStep.ProvisioningStepToPlan) return;
    const stepPlan = this.envService.getPlan(this.provisioningStep.ProvisioningStepToPlan.id);
    if (!stepPlan) return;
    this.latestDiff = [...stepPlan.PlanToPlanDiffs].sort((a, b) => b.revision - a.revision)[0];
  }

  getStatusIcon(): string {
    switch (this.provisioningStep.type) {
      case LaForgeProvisioningStepType.Script:
        return 'file-code';
      case LaForgeProvisioningStepType.Command:
        return 'terminal';
      case LaForgeProvisioningStepType.DnsRecord:
        return 'globe';
      case LaForgeProvisioningStepType.FileDownload:
        return 'download';
      case LaForgeProvisioningStepType.FileDelete:
        return 'trash';
      case LaForgeProvisioningStepType.FileExtract:
        return 'file-archive';
      default:
        return 'minus-circle';
    }
  }

  getStatusColor(): string {
    if (this.mode === 'plan') {
      // if (!this.latestDiff) return 'dark';
      // switch (this.latestDiff.new_state) {
      //   case LaForgeProvisionStatus.Torebuild:
      //     return 'warning';
      //   default:
      //     return 'dark';
      // }
      return 'black';
    }
    // const status = this.provisionedStep.ProvisioningStepToPlan?.PlanToStatus ?? this.provisionedStep.ProvisioningStepToStatus;
    const status = this.planStatus ?? this.provisioningStepStatus;
    if (!status?.state) {
      return 'black';
    }
    switch (status.state) {
      case LaForgeProvisionStatus.Complete:
        return 'success';
      case LaForgeProvisionStatus.Todelete:
        return 'warning';
      case LaForgeProvisionStatus.Deleteinprogress:
        return 'warning';
      case LaForgeProvisionStatus.Deleted:
        return 'dark';
      case LaForgeProvisionStatus.Failed:
        return 'danger';
      case LaForgeProvisionStatus.Inprogress:
        return 'info';
      case LaForgeProvisionStatus.Planning:
        return 'primary';
      default:
        return 'dark';
    }
  }

  getText(): string {
    switch (this.provisioningStep.type) {
      case LaForgeProvisioningStepType.Script:
        return `${this.provisioningStep.ProvisioningStepToScript.source} ${this.provisioningStep.ProvisioningStepToScript.args.join(' ')}`;
      case LaForgeProvisioningStepType.Command:
        return `${this.provisioningStep.ProvisioningStepToCommand.program} ${this.provisioningStep.ProvisioningStepToCommand.args.join(
          ' '
        )}`;
      case LaForgeProvisioningStepType.DnsRecord:
        return 'DNSRecord';
      case LaForgeProvisioningStepType.FileDownload:
        // eslint-disable-next-line max-len
        return `${this.provisioningStep.ProvisioningStepToFileDownload.source} -> ${this.provisioningStep.ProvisioningStepToFileDownload.destination}`;
      case LaForgeProvisioningStepType.FileDelete:
        return 'FileDelete';
      case LaForgeProvisioningStepType.FileExtract:
        return 'FileExtract';
      default:
        return 'Step';
    }
  }
}
