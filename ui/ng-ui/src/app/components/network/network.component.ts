import { Component, Input, OnInit, ChangeDetectorRef, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { LaForgeProvisionedNetwork, LaForgeProvisionStatus, LaForgeSubscribeUpdatedStatusSubscription } from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';
import { Subscription } from 'rxjs';
import { Status } from 'src/app/models/common.model';
import { RebuildService } from 'src/app/services/rebuild/rebuild.service';

import { NetworkModalComponent } from '../network-modal/network-modal.component';

@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss']
})
export class NetworkComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
  @Input() provisionedNetwork: LaForgeProvisionedNetwork;
  @Input() status: Status;
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() parentSelected: boolean;
  @Input() mode: 'plan' | 'build' | 'manage';
  isSelectedState = false;
  planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  expandOverride = false;

  constructor(
    public dialog: MatDialog,
    private rebuild: RebuildService,
    private envService: EnvironmentService,
    private cdRef: ChangeDetectorRef
  ) {
    if (!this.mode) this.mode = 'manage';
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
    if (!this.parentSelected) this.parentSelected = false;
  }

  ngOnInit(): void {
    const sub1 = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkPlanStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub1);
    const sub2 = this.envService.agentStatusUpdate.asObservable().subscribe(() => {
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub2);
  }

  ngOnDestroy() {
    this.unsubscribe.forEach((s) => s.unsubscribe());
  }

  checkPlanStatus(): void {
    this.planStatus = this.envService.getStatus(this.provisionedNetwork.ProvisionedNetworkToPlan.PlanToStatus.id) || this.planStatus;
  }

  viewDetails(): void {
    this.dialog.open(NetworkModalComponent, {
      width: '50%',
      height: '80%',
      data: { provisionedNetwork: this.provisionedNetwork, planStatus: this.planStatus }
    });
  }

  allChildrenResponding(): boolean {
    let numWithAgentData = 0;
    let numWithCompletedSteps = 0;
    let totalHosts = 0;
    for (const host of this.provisionedNetwork.ProvisionedNetworkToProvisionedHost) {
      totalHosts++;
      if (host.ProvisionedHostToAgentStatus?.clientId) numWithAgentData++;
      let totalSteps = 0;
      let totalCompletedSteps = 0;
      for (const step of host.ProvisionedHostToProvisioningStep) {
        if (step.step_number === 0) continue;
        totalSteps++;
        if (
          step.ProvisioningStepToStatus.id &&
          this.envService.getStatus(step.ProvisioningStepToPlan.PlanToStatus.id)?.state === LaForgeProvisionStatus.Complete
        )
          totalCompletedSteps++;
      }
      if (totalSteps === totalCompletedSteps) numWithCompletedSteps++;
    }
    return numWithAgentData === totalHosts && numWithCompletedSteps === totalHosts;
  }

  getStatusIcon(): string {
    if (!this.planStatus) return 'minus-circle';

    switch (this.planStatus.state) {
      case LaForgeProvisionStatus.Planning:
        return 'ruler-triangle fas';
      case LaForgeProvisionStatus.Todelete:
        return 'recycle fas';
      case LaForgeProvisionStatus.Deleteinprogress:
        return 'trash-restore fas';
      case LaForgeProvisionStatus.Deleted:
        return 'trash fas';
      default:
        return 'network-wired';
    }
  }

  getStatusColor(): string {
    if (!this.planStatus) return 'dark';
    switch (this.planStatus.state) {
      case LaForgeProvisionStatus.Complete:
        if (this.allChildrenResponding()) {
          return 'success';
        } else {
          return 'warning';
        }
      case LaForgeProvisionStatus.Inprogress:
        return 'info';
      case LaForgeProvisionStatus.Tainted:
        return 'danger';
      case LaForgeProvisionStatus.Failed:
        return 'danger';
      case LaForgeProvisionStatus.Todelete:
        return 'primary';
      case LaForgeProvisionStatus.Deleteinprogress:
        return 'info';
      case LaForgeProvisionStatus.Planning:
        return 'primary';
      default:
        return 'dark';
    }
  }

  onSelect(): void {
    let success = false;
    if (!this.isSelected()) {
      success = this.rebuild.addNetwork(this.provisionedNetwork);
    } else {
      success = this.rebuild.removeNetwork(this.provisionedNetwork);
    }
    if (success) this.isSelectedState = !this.isSelectedState;
  }

  onIndeterminateChange(isIndeterminate: boolean): void {
    if (!isIndeterminate && this.isSelectedState) setTimeout(() => this.rebuild.addNetwork(this.provisionedNetwork), 500);
  }

  isSelected(): boolean {
    return this.rebuild.hasNetwork(this.provisionedNetwork);
  }

  shouldCollapse(): boolean {
    if (this.mode === 'plan') {
      const plan = this.envService.getPlan(this.provisionedNetwork.ProvisionedNetworkToPlan.id);
      if (plan?.PlanToPlanDiffs.length > 0) {
        const latestDiff = plan.PlanToPlanDiffs.sort((a, b) => b.revision - a.revision)[0];
        if (latestDiff.new_state === LaForgeProvisionStatus.Planning) {
          return false;
        } else {
          return true;
        }
      }
      return false;
    }
    return (
      this.planStatus &&
      (this.planStatus.state === LaForgeProvisionStatus.Deleted ||
        (this.planStatus.state === LaForgeProvisionStatus.Complete && this.allChildrenResponding()))
    );
  }

  canOverrideExpand(): boolean {
    return (
      this.planStatus &&
      (this.planStatus.state === LaForgeProvisionStatus.Complete || this.planStatus.state === LaForgeProvisionStatus.Deleted)
    );
  }

  toggleCollapse(): void {
    this.expandOverride = !this.expandOverride;
  }
}
