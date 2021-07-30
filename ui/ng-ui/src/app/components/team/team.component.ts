import { Component, Input, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import {
  LaForgeProvisionStatus,
  LaForgeSubscribeUpdatedStatusSubscription,
  LaForgeTeam,
  LaForgeGetBuildTreeQuery,
  LaForgePlanFieldsFragment
} from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';
import { Subscription } from 'rxjs';

import { RebuildService } from '../../services/rebuild/rebuild.service';

@Component({
  selector: 'app-team',
  templateUrl: './team.component.html',
  styleUrls: ['./team.component.scss']
})
export class TeamComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
  @Input() title: string;
  @Input() team: LaForgeGetBuildTreeQuery['build']['buildToTeam'][0];
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() mode: 'plan' | 'build' | 'manage';
  isSelectedState = false;
  planStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  expandOverride = false;
  shouldHideLoading = false;
  shouldHide = false;
  latestDiff: LaForgePlanFieldsFragment['PlanToPlanDiffs'][0];

  constructor(private rebuild: RebuildService, private envService: EnvironmentService, private cdRef: ChangeDetectorRef) {
    if (!this.mode) this.mode = 'manage';
    if (!this.style) this.style = 'compact';
    if (!this.selectable) this.selectable = false;
  }

  ngOnInit(): void {
    const sub1 = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkPlanStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub1);
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

  checkPlanStatus(): void {
    this.planStatus = this.envService.getStatus(this.team.TeamToPlan.PlanToStatus.id) || this.planStatus;
  }

  checkLatestPlanDiff(): void {
    if (this.latestDiff) return;
    const teamPlan = this.envService.getPlan(this.team.TeamToPlan.id);
    if (!teamPlan) return;
    this.latestDiff = [...teamPlan.PlanToPlanDiffs].sort((a, b) => b.revision - a.revision)[0];
  }

  // getStatus(): ProvisionStatus {
  //   // let status: ProvisionStatus = ProvisionStatus.ProvStatusUndefined;
  //   let numWithAgentData = 0;
  //   let totalAgents = 0;
  //   for (const network of this.team.TeamToProvisionedNetwork) {
  //     for (const host of network.ProvisionedNetworkToProvisionedHost) {
  //       totalAgents++;
  //       if (host.ProvisionedHostToAgentStatus?.clientId) numWithAgentData++;
  //     }
  //   }
  //   if (numWithAgentData === totalAgents) return ProvisionStatus.COMPLETE;
  //   else if (numWithAgentData === 0) return ProvisionStatus.FAILED;
  //   else return ProvisionStatus.INPROGRESS;
  // }

  allChildrenResponding(): boolean {
    let numWithAgentData = 0;
    let numWithCompletedSteps = 0;
    let totalHosts = 0;
    for (const pnet of this.team.TeamToProvisionedNetwork) {
      for (const host of pnet.ProvisionedNetworkToProvisionedHost) {
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
    }
    return numWithAgentData === totalHosts && numWithCompletedSteps === totalHosts;
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
          return 'fal fa-users';
      }
    }
    if (!this.planStatus) return 'fas fa-minus-circle';

    switch (this.planStatus.state) {
      case LaForgeProvisionStatus.Planning:
        return 'fas fa-ruler-triangle';
      case LaForgeProvisionStatus.Todelete:
        return 'fas fa-recycle fas';
      case LaForgeProvisionStatus.Deleteinprogress:
        return 'fas fa-trash-restore';
      case LaForgeProvisionStatus.Deleted:
        return 'fas fa-trash fas';
      case LaForgeProvisionStatus.Failed:
        return 'fas fa-ban';
      default:
        return 'fal fa-users';
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

  // getStatusColor(): string {
  //   switch (this.getStatus()) {
  //     case ProvisionStatus.COMPLETE:
  //       return 'success';
  //     case ProvisionStatus.INPROGRESS:
  //       return 'warning';
  //     case ProvisionStatus.FAILED:
  //       return 'danger';
  //     default:
  //       return 'dark';
  //   }
  // }

  onSelect(): void {
    let success = false;
    if (!this.isSelected()) {
      success = this.rebuild.addTeam(this.team as LaForgeTeam);
    } else {
      success = this.rebuild.removeTeam(this.team as LaForgeTeam);
    }
    if (success) this.isSelectedState = !this.isSelectedState;
  }

  isSelected(): boolean {
    return this.rebuild.hasTeam(this.team as LaForgeTeam);
  }

  checkShouldHide() {
    if (this.mode === 'plan') {
      if (!this.latestDiff) return (this.shouldHide = false);
      const latestCommit = this.envService.getLatestCommit();
      if (!latestCommit) return false;
      const teamPlan = this.envService.getPlan(this.team.TeamToPlan.id);
      if (teamPlan?.PlanToPlanDiffs.length > 0) {
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
      //   const latestCommit = this.envService.getBuildTree().getValue()?.BuildToLatestBuildCommit;
      //   const teamPlan = this.envService.getPlan(this.team.TeamToPlan.id);
      //   if (teamPlan?.PlanToPlanDiffs.length > 0) {
      //     const latestDiff = [...teamPlan.PlanToPlanDiffs].sort((a, b) => b.revision - a.revision)[0];
      //     // expand if latest diff is a part of the latest commit
      //     if (latestCommit && latestCommit.BuildCommitToPlanDiffs.filter((diff) => diff.id === latestDiff.id).length > 0) {
      //       // this.expandOverride = true;
      return false;
      //     }
      //   }
      //   // this.expandOverride = false;
      //   return true;
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
