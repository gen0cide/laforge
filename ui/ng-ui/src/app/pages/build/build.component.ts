import { ChangeDetectorRef, Component, OnInit, OnDestroy } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import {
  LaForgeGetBuildTreeQuery,
  LaForgeGetEnvironmentInfoQuery,
  LaForgeProvisionStatus,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { Observable, Subscription } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

import { LaForgeExecuteBuildGQL } from '../../../generated/graphql';

@Component({
  selector: 'app-build',
  templateUrl: './build.component.html',
  styleUrls: ['./build.component.scss']
})
export class BuildComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
  environment: Observable<LaForgeGetEnvironmentInfoQuery['environment']>;
  build: Observable<LaForgeGetBuildTreeQuery['build']>;
  buildStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  executeBuildLoading = false;
  planStatusesLoading = false;
  agentStatusesLoading = false;

  constructor(
    private subheader: SubheaderService,
    public envService: EnvironmentService,
    private cdRef: ChangeDetectorRef,
    private executeBuild: LaForgeExecuteBuildGQL,
    private snackbar: MatSnackBar
  ) {
    this.subheader.setTitle('Build');
    this.subheader.setDescription('Monitor the progress of a given build');
    this.subheader.setShowEnvDropdown(true);

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
    this.envService.buildIsLoading.subscribe((isLoading) => {
      if (isLoading)
        this.snackbar.open('Environment is loading...', null, {
          panelClass: ['bg-info', 'text-white']
        });
      else if (!this.envService.buildIsLoading.getValue()) this.snackbar.dismiss();
    });
    this.envService.buildIsLoading.subscribe((isLoading) => {
      if (isLoading)
        this.snackbar.open('Build is loading...', null, {
          panelClass: ['bg-info', 'text-white']
        });
      else if (!this.envService.envIsLoading.getValue()) this.snackbar.dismiss();
    });
  }

  ngOnInit(): void {
    const sub1 = this.envService.getBuildTree().subscribe((buildTree) => {
      if (!buildTree) return;
      this.planStatusesLoading = true;
      this.envService
        .initPlanStatuses()
        .catch((err) => {
          this.snackbar.open(err, 'Okay', {
            panelClass: ['bg-danger', 'text-white']
          });
        })
        .finally(() => {
          this.planStatusesLoading = false;
          this.cdRef.detectChanges();
        });
      this.agentStatusesLoading = true;
      this.envService
        .initAgentStatuses()
        .catch((err) => {
          this.snackbar.open(err, 'Okay', {
            panelClass: ['bg-danger', 'text-white']
          });
        })
        .finally(() => {
          this.agentStatusesLoading = false;
          this.cdRef.detectChanges();
        });
      this.envService.startAgentStatusSubscription();
      this.envService.startStatusSubscription();
      this.envService.startAgentTaskSubscription();
      this.envService.startBuildCommitSubscription();
      this.envService.startBuildSubscription();
    });
    this.unsubscribe.push(sub1);
    const sub2 = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkBuildStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub2);
  }

  ngOnDestroy(): void {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
    this.envService.stopAgentStatusSubscription();
    this.envService.stopStatusSubscription();
    this.envService.stopAgentTaskSubscription();
    this.envService.stopBuildCommitSubscription();
    this.envService.stopBuildSubscription();
  }

  checkBuildStatus(): void {
    if (!this.envService.getBuildTree().getValue()) return;
    const updatedStatus = this.envService.getStatus(this.envService.getBuildTree().getValue().buildToStatus.id);
    if (updatedStatus) {
      this.buildStatus = { ...updatedStatus };
    }
  }

  envIsSelected(): boolean {
    return this.envService.getEnvironmentInfo().getValue() != null;
  }

  triggerExecuteBuild(): void {
    if (!this.envService.getBuildTree().getValue()?.id) return;
    this.executeBuildLoading = true;
    this.executeBuild
      .mutate({
        buildId: this.envService.getBuildTree().getValue().id
      })
      .toPromise()
      .then(({ data, errors }) => {
        if (errors) {
          return console.error(errors);
        }
      }, console.error)
      .finally(() => {
        this.executeBuildLoading = false;
        this.snackbar.open('Successfully started build!', 'Cool', {
          duration: 3000
        });
      });
  }

  canExecuteBuild(): boolean {
    return this.buildStatus && this.buildStatus.state === LaForgeProvisionStatus.Planning;
  }
}
