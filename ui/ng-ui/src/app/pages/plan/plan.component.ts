import { ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import {
  LaForgeApproveBuildCommitGQL,
  LaForgeBuildCommitFieldsFragment,
  LaForgeBuildCommitState,
  LaForgeCancelBuildCommitGQL,
  LaForgeGetBuildTreeQuery,
  LaForgeGetEnvironmentInfoQuery,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
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
  latestCommit: BehaviorSubject<LaForgeBuildCommitFieldsFragment>;
  apolloError: any = {};
  approveDenyCommitLoading = false;
  planLoading = true;

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    public envService: EnvironmentService,
    private approveBuildCommit: LaForgeApproveBuildCommitGQL,
    private cancelBuildCommit: LaForgeCancelBuildCommitGQL,
    private snackBar: MatSnackBar,
    private router: Router
  ) {
    this.subheader.setTitle('Plan');
    this.subheader.setDescription('Plan an environment to build');
    this.subheader.setShowEnvDropdown(true);

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
    this.latestCommit = new BehaviorSubject(null);
  }

  ngOnInit(): void {
    this.planLoading = true;
    const sub1 = this.envService.getBuildTree().subscribe(() => {
      this.envService.initPlanStatuses();
      // this.envService.initAgentStatuses();
      this.envService.initPlans();
      this.envService.initBuildCommits();
      this.checkLatestCommit();
      this.envService.startStatusSubscription();
      this.envService.startAgentStatusSubscription();
      this.envService.startBuildSubscription();
      this.envService.startBuildCommitSubscription();
    });
    this.unsubscribe.push(sub1);
    const sub2 = this.envService.statusUpdate.subscribe(() => {
      this.checkBuildStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub2);
    this.planLoading = true;
    const sub3 = this.envService.planUpdate.subscribe(() => {
      this.planLoading = false;
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub3);
    // const sub4 = this.envService.buildCommitUpdate.subscribe(() => {
    //   this.checkLatestCommit();
    //   this.cdRef.detectChanges();
    // });
    // this.unsubscribe.push(sub4);
  }

  ngOnDestroy(): void {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
    this.envService.stopStatusSubscription();
    this.envService.stopAgentStatusSubscription();
    this.envService.stopBuildSubscription();
    this.envService.stopBuildCommitSubscription();
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

  checkLatestCommit(): void {
    if (!this.envService.getBuildTree().getValue()) return;
    const updatedBuildCommit = this.envService.getLatestCommit();
    if (updatedBuildCommit && (!this.latestCommit.getValue() || this.latestCommit.getValue().id !== updatedBuildCommit.id)) {
      this.latestCommit.next({ ...updatedBuildCommit });
      this.cdRef.detectChanges();
    }
    if (this.latestCommit.getValue() && this.latestCommit.getValue().state === LaForgeBuildCommitState.Inprogress) {
      this.router.navigate(['build']);
      this.snackBar.open('Plan is in progress. Redirecting to build...', null, {
        duration: 1000
      });
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
    const latestCommit = this.envService.getLatestCommit();
    if (!latestCommit) return '';
    switch (latestCommit.state) {
      case LaForgeBuildCommitState.Approved:
        return 'accent';
      case LaForgeBuildCommitState.Cancelled:
        return 'warn';
      case LaForgeBuildCommitState.Inprogress:
        return 'accent';
      case LaForgeBuildCommitState.Planning:
        return 'primary';
      case LaForgeBuildCommitState.Applied:
        return 'link';
    }
  }

  getCommitStateText(): string {
    const latestCommit = this.envService.getLatestCommit();
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
    this.approveBuildCommit
      .mutate({
        buildCommitId: this.envService.getBuildTree().getValue().BuildToLatestBuildCommit.id
      })
      .toPromise()
      .then(
        ({ data, errors }) => {
          if (errors) {
            this.snackBar.open('Error while approving commit. See logs for more info.', 'Okay', {
              duration: 3000,
              panelClass: 'bg-danger text-white'
            });
          } else if (data.approveCommit) {
            this.snackBar.open('Commit approved', 'Nice!', {
              // duration: 3000,
              panelClass: ['bg-success', 'text-white']
            });
            this.router.navigate(['build']);
          }
        },
        (err) => {
          console.error(err);
          this.snackBar.open('Error while approving commit. See console for more info.', 'Okay', {
            duration: 3000,
            panelClass: 'bg-danger text-white'
          });
        }
      );
  }

  cancelCommit(): void {
    this.cancelBuildCommit
      .mutate({
        buildCommitId: this.envService.getBuildTree().getValue().BuildToLatestBuildCommit.id
      })
      .toPromise()
      .then(
        ({ data, errors }) => {
          if (errors) {
            this.snackBar.open('Error while cancelling commit. See logs for more info.', 'Okay', {
              duration: 3000,
              panelClass: 'bg-danger'
            });
          } else if (data.cancelCommit) {
            this.snackBar.open('Commit cancelled', 'Okay', {
              duration: 3000
            });
          }
        },
        (err) => {
          console.error(err);
          this.snackBar.open('Error while cancelling commit. See console for more info.', 'Okay', {
            duration: 3000,
            panelClass: 'bg-danger'
          });
        }
      );
  }

  canApproveDenyCommit(): boolean {
    const latestCommit = this.envService.getLatestCommit();
    if (!latestCommit) return false;
    if (latestCommit.state === LaForgeBuildCommitState.Planning) return true;
    return false;
  }
}
