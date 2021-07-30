import { Component, Input, OnInit, ChangeDetectorRef, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {
  LaForgeGetBuildTreeQuery,
  LaForgePlanFieldsFragment,
  LaForgeProvisionedHost,
  LaForgeProvisionStatus,
  LaForgeSubscribeUpdatedAgentStatusSubscription,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';

import { Subscription } from 'rxjs';
import { ApiService } from 'src/app/services/api/api.service';

import { RebuildService } from '../../services/rebuild/rebuild.service';
import { HostModalComponent } from '../host-modal/host-modal.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.scss']
})
export class HostComponent implements OnInit, OnDestroy {
  // @Input() status: Status;
  @Input()
  // eslint-disable-next-line max-len
  provisionedHost: LaForgeGetBuildTreeQuery['build']['buildToTeam'][0]['TeamToProvisionedNetwork'][0]['ProvisionedNetworkToProvisionedHost'][0];
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() parentSelected: boolean;
  @Input() hasAgent: boolean;
  @Input() mode: 'plan' | 'build' | 'manage';
  unsubscribe: Subscription[] = [];
  isSelectedState = false;
  planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  provisionedHostStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  agentStatus: LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'];
  expandOverride = false;
  shouldHideLoading = false;
  shouldHide = false;
  latestDiff: LaForgePlanFieldsFragment['PlanToPlanDiffs'][0];

  constructor(
    public dialog: MatDialog,
    private rebuild: RebuildService,
    private api: ApiService,
    private envService: EnvironmentService,
    private cdRef: ChangeDetectorRef
  ) {
    if (!this.mode) this.mode = 'manage';
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
    if (!this.parentSelected) this.parentSelected = false;
    if (!this.hasAgent) this.hasAgent = false;
  }

  ngOnInit() {
    const sub1 = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkPlanStatus();
      this.checkProvisionedHostStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub1);
    const sub2 = this.envService.agentStatusUpdate.asObservable().subscribe(() => {
      this.checkAgentStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub2);
    if (this.mode === 'plan') {
      this.shouldHideLoading = true;
      const sub2 = this.envService.planUpdate.asObservable().subscribe(() => {
        this.checkLatestPlanDiff();
        this.checkShouldHide();
        this.cdRef.detectChanges();
      });
      this.unsubscribe.push(sub2);
      const sub4 = this.envService.buildCommitUpdate.asObservable().subscribe(() => {
        this.checkLatestPlanDiff();
        this.checkShouldHide();
        this.cdRef.detectChanges();
      });
      this.unsubscribe.push(sub4);
    }
  }

  ngOnDestroy() {
    this.unsubscribe.forEach((s) => s.unsubscribe());
  }

  viewDetails(): void {
    this.dialog.open(HostModalComponent, {
      width: '50%',
      height: '80%',
      data: {
        provisionedHost: this.provisionedHost,
        planStatus: this.planStatus,
        provisionedHostStatus: this.provisionedHostStatus,
        agentStatus: this.agentStatus
      }
    });
  }

  checkPlanStatus(): void {
    this.planStatus = this.envService.getStatus(this.provisionedHost.ProvisionedHostToPlan.PlanToStatus.id) || this.planStatus;
  }

  checkProvisionedHostStatus(): void {
    this.provisionedHostStatus = this.envService.getStatus(this.provisionedHost.ProvisionedHostToStatus.id) || this.provisionedHostStatus;
  }

  checkAgentStatus(): void {
    // if (this.provisionedHost.ProvisionedHostToAgentStatus?.clientId) {
    const updatedStatus = this.envService.getAgentStatus(this.provisionedHost.id);
    if (updatedStatus) {
      // console.log('agent status updated');
      this.agentStatus = updatedStatus;
    }
    // }
  }

  isAgentStale(): boolean {
    if (!this.agentStatus) return true;
    return Date.now() / 1000 - this.agentStatus.timestamp > 120;
  }

  checkLatestPlanDiff(): void {
    if (this.latestDiff) return;
    const phostPlan = this.envService.getPlan(this.provisionedHost.ProvisionedHostToPlan.id);
    if (!phostPlan) return;
    this.latestDiff = [...phostPlan.PlanToPlanDiffs].sort((a, b) => b.revision - a.revision)[0];
  }

  getStatusIcon(): string {
    if (this.mode === 'plan') {
      if (!this.latestDiff) return 'fas fa-spinner fa-spin';
      switch (this.latestDiff.new_state) {
        case LaForgeProvisionStatus.Torebuild:
          return 'fas fa-sync-alt';
        case LaForgeProvisionStatus.Todelete:
          return 'fad fa-trash';
        case LaForgeProvisionStatus.Planning:
          return 'fas fa-ruler-triangle';
        default:
          return 'fas fa-computer-classic';
      }
    }
    const status = this.planStatus ?? this.provisionedHostStatus;
    if (status?.state) {
      switch (status.state) {
        case LaForgeProvisionStatus.Todelete:
          return 'fas fa-recycle';
        case LaForgeProvisionStatus.Deleteinprogress:
          return 'fas fa-trash-restore';
        case LaForgeProvisionStatus.Deleted:
          return 'fad fa-trash';
      }
    }
    if (this.agentStatus) {
      if (this.isAgentStale()) return 'fas fa-exclamation-circle';
      if (this.childrenCompleted()) return 'fas fa-check-circle';
      else return 'fas fa-satellite-dish';
    } else {
      if (!status?.state) {
        return 'fas fa-minus-circle';
      }
      switch (status.state) {
        case LaForgeProvisionStatus.Complete:
          return 'fas fa-box-check';
        case LaForgeProvisionStatus.Failed:
          return 'fas fa-ban';
        case LaForgeProvisionStatus.Inprogress:
          return 'fas fa-play-circle';
        case LaForgeProvisionStatus.Awaiting:
          return 'fas fa-spinner fa-spin';
        case LaForgeProvisionStatus.Planning:
          return 'fas fa-ruler-triangle';
        default:
          return 'fas fa-computer-classic';
      }
    }
  }

  getStatusColor(): string {
    if (this.mode === 'plan') {
      if (!this.latestDiff) return 'dark';
      switch (this.latestDiff.new_state) {
        case LaForgeProvisionStatus.Torebuild:
          return 'warning';
        case LaForgeProvisionStatus.Todelete:
          return 'danger';
        case LaForgeProvisionStatus.Planning:
          return 'primary';
        default:
          return 'dark';
      }
    }
    const status = this.planStatus ?? this.provisionedHostStatus;
    if (status?.state) {
      switch (status.state) {
        case LaForgeProvisionStatus.Todelete:
          return 'primary';
        case LaForgeProvisionStatus.Deleteinprogress:
          return 'info';
        case LaForgeProvisionStatus.Deleted:
          return 'dark';
      }
    }
    if (this.agentStatus) {
      if (this.isAgentStale()) return 'warning';
      else return 'success';
    } else {
      if (!status?.state) {
        return 'minus-circle';
      }
      switch (status.state) {
        case LaForgeProvisionStatus.Complete:
          return 'success';
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
  }

  onSelect(): void {
    let success = false;
    if (!this.isSelected()) {
      success = this.rebuild.addHost(this.provisionedHost as LaForgeProvisionedHost);
    } else {
      success = this.rebuild.removeHost(this.provisionedHost as LaForgeProvisionedHost);
    }
    // console.log(success);
    if (success) this.isSelectedState = !this.isSelectedState;
  }

  onIndeterminateChange(isIndeterminate: boolean): void {
    if (!isIndeterminate && this.isSelectedState)
      setTimeout(() => this.rebuild.addHost(this.provisionedHost as LaForgeProvisionedHost), 500);
  }

  isSelected(): boolean {
    return this.rebuild.hasHost(this.provisionedHost as LaForgeProvisionedHost);
  }

  childrenCompleted(): boolean {
    let numCompleted = 0;
    let numAwaiting = 0;
    let totalSteps = 0;
    for (const step of this.provisionedHost.ProvisionedHostToProvisioningStep) {
      if (step.step_number === 0) continue;
      totalSteps++;
      const stepStatus = this.envService.getStatus(step.ProvisioningStepToPlan.PlanToStatus.id);
      if (stepStatus?.state === LaForgeProvisionStatus.Complete) numCompleted++;
      if (stepStatus?.state === LaForgeProvisionStatus.Awaiting) numAwaiting++;
      // if (step.ProvisioningStepToStatus.state === LaForgeProvisionStatus.Complete) numCompleted++;
    }
    // console.log(this.provisionedHost.ProvisionedHostToHost.hostname, totalSteps, numCompleted);
    if (numCompleted === totalSteps) return true;
    if (numAwaiting === totalSteps) return true;
    else return false;
  }

  checkShouldHide() {
    if (this.mode === 'plan') {
      if (!this.latestDiff) return (this.shouldHide = false);
      const latestCommit = this.envService.getLatestCommit();
      if (!latestCommit) return (this.shouldHide = false);
      const phostPlan = this.envService.getPlan(this.provisionedHost.ProvisionedHostToPlan.id);
      if (phostPlan?.PlanToPlanDiffs.length > 0) {
        // expand if latest diff is a part of the latest commit
        if (latestCommit && latestCommit.BuildCommitToPlanDiffs.filter((diff) => diff.id === this.latestDiff.id).length > 0) {
          this.shouldHideLoading = false;
          this.shouldHide = false;
          return;
        }
      }
      this.shouldHideLoading = false;
      this.shouldHide = true;
      return;
    }
    this.shouldHide = false;
  }

  shouldCollapse(): boolean {
    if (this.mode === 'plan') {
      // const plan = this.envService.getPlan(this.provisionedHost.ProvisionedHostToPlan.id);
      // if (plan?.PlanToPlanDiffs.length > 0) {
      //   const latestCommitRevision = this.envService.getBuildTree().getValue()?.BuildToLatestBuildCommit.revision;
      //   const latestDiff = [...plan.PlanToPlanDiffs].sort((a, b) => b.revision - a.revision)[0];
      //   // collapse if latest diff isn't a part of the latest commit
      //   if (latestCommitRevision && latestCommitRevision != latestDiff.revision) return true;
      // return false;
      // }
      return true;
    }
    if (this.planStatus && this.planStatus.state === LaForgeProvisionStatus.Deleted) return true;
    // return hostChildrenCompleted(this.provisionedHost as LaForgeProvisionedHost, this.envService.getStatus);
    return this.childrenCompleted();
  }

  toggleCollapse(): void {
    this.expandOverride = !this.expandOverride;
  }
}
