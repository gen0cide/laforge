import { ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import {
  LaForgeBuildCommitState,
  LaForgeExecuteBuildGQL,
  LaForgeGetBuildTreeQuery,
  LaForgeGetEnvironmentInfoQuery,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { Observable, Subscription } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

@Component({
  selector: 'app-plan',
  templateUrl: './plan.component.html',
  styleUrls: ['./plan.component.scss']
})
export class PlanComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
  environment: Observable<LaForgeGetEnvironmentInfoQuery['environment']>;
  build: Observable<LaForgeGetBuildTreeQuery['build']>;
  buildStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  apolloError: any = {};
  approveDenyCommitLoading = false;
  planLoading = true;

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    public envService: EnvironmentService,
    private executeBuild: LaForgeExecuteBuildGQL,
    private snackBar: MatSnackBar,
    private router: Router
  ) {
    this.subheader.setTitle('Plan');
    this.subheader.setDescription('Plan an environment to build');
    this.subheader.setShowEnvDropdown(true);

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
    this.planLoading = true;
  }

  ngOnInit(): void {
    const sub1 = this.envService.getBuildTree().subscribe(() => {
      this.envService.initPlanStatuses();
      this.envService.initAgentStatuses();
      this.envService.initPlans();
    });
    this.unsubscribe.push(sub1);
    const sub2 = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkBuildStatus();
      this.cdRef.markForCheck();
    });
    this.unsubscribe.push(sub2);
    this.planLoading = true;
    const sub3 = this.envService.planUpdate.asObservable().subscribe(() => {
      this.planLoading = false;
      // this.cdRef.markForCheck();
    });
    this.unsubscribe.push(sub3);
  }

  ngOnDestroy(): void {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
  }

  envIsSelected(): boolean {
    return this.envService.getEnvironmentInfo().getValue() != null;
  }

  checkBuildStatus(): void {
    if (!this.envService.getBuildTree().getValue()) return;
    const updatedStatus = this.envService.getStatus(this.envService.getBuildTree().getValue().buildToStatus.id);
    if (updatedStatus) {
      this.buildStatus = { ...updatedStatus };
    }
  }

  // triggerExecuteBuild(): void {
  //   if (!this.envService.getBuildTree().getValue()?.id) return;
  //   this.executeBuildLoading = true;
  //   this.executeBuild
  //     .mutate({
  //       buildId: this.envService.getBuildTree().getValue().id
  //     })
  //     .toPromise()
  //     .then(({ data, errors }) => {
  //       if (errors) {
  //         return console.error(errors);
  //       } else {
  //         this.snackBar.open('Successfully started build!', 'Cool', {
  //           duration: 3000
  //         });
  //         this.router.navigate(['build']);
  //       }
  //     }, console.error)
  //     .finally(() => {
  //       this.executeBuildLoading = false;
  //     });
  // }

  // canExecuteBuild(): boolean {
  //   return this.buildStatus && this.buildStatus.state === LaForgeProvisionStatus.Planning;
  // }

  getCommitStateColor(): string {
    const latestCommit = this.envService.getBuildTree().getValue()?.BuildToLatestBuildCommit;
    if (!latestCommit) return '';
    switch (latestCommit.state) {
      case LaForgeBuildCommitState.Approved:
        return 'primary';
      case LaForgeBuildCommitState.Cancelled:
        return 'warn';
      case LaForgeBuildCommitState.Inprogress:
        return 'accent';
      case LaForgeBuildCommitState.Planning:
        return '';
      case LaForgeBuildCommitState.Applied:
        return 'primary';
    }
  }

  getCommitStateText(): string {
    const latestCommit = this.envService.getBuildTree().getValue()?.BuildToLatestBuildCommit;
    if (!latestCommit) return '';
    switch (latestCommit.state) {
      case LaForgeBuildCommitState.Approved:
        return 'Approved';
      case LaForgeBuildCommitState.Cancelled:
        return 'Cancelled';
      case LaForgeBuildCommitState.Inprogress:
        return 'In Progress';
      case LaForgeBuildCommitState.Planning:
        return 'Planning';
      case LaForgeBuildCommitState.Applied:
        return 'Applied';
    }
  }

  approveCommit(): void {
    console.log('TODO: Add commit approval');
  }

  cancelCommit(): void {
    console.log('TODO: Add commit cancellation');
  }

  canApproveDenyCommit(): boolean {
    const latestCommit = this.envService.getBuildTree().getValue()?.BuildToLatestBuildCommit;
    if (!latestCommit) return false;
    if (latestCommit.state === LaForgeBuildCommitState.Planning) return true;
    return false;
  }
}
