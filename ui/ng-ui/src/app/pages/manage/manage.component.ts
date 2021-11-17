import { ChangeDetectorRef, Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import {
  LaForgeGetEnvironmentInfoQuery,
  LaForgeGetBuildTreeQuery,
  LaForgeSubscribeUpdatedStatusSubscription,
  LaForgeProvisionStatus
} from '@graphql';
import { RebuildService } from '@services/rebuild/rebuild.service';
import { GraphQLError } from 'graphql';
import { Observable, Subscription } from 'rxjs';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

import { SubheaderService } from '../../_metronic/partials/layout/subheader/_services/subheader.service';

import { DeleteBuildModalComponent } from '@components/delete-build-modal/delete-build-modal.component';

@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ManageComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
  environment: Observable<LaForgeGetEnvironmentInfoQuery['environment']>;
  build: Observable<LaForgeGetBuildTreeQuery['build']>;
  envIsLoading: Observable<boolean>;
  buildIsLoading: Observable<boolean>;
  environmentDetailsCols: string[] = ['TeamCount', 'AdminCIDRs', 'ExposedVDIPorts'];
  selectionMode = false;
  isRebuildLoading = false;
  rebuildErrors: (GraphQLError | Error)[] = [];
  confirmDeleteBuild = false;
  buildStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];
  planStatusesLoading = false;
  agentStatusesLoading = false;

  constructor(
    private dialog: MatDialog,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    private envService: EnvironmentService,
    private rebuild: RebuildService,
    private router: Router,
    private snackbar: MatSnackBar
  ) {
    this.subheader.setTitle('Environment');
    this.subheader.setDescription('Manage your currently running environment');
    this.subheader.setShowEnvDropdown(true);

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
    this.envIsLoading = this.envService.envIsLoading.asObservable();
    this.buildIsLoading = this.envService.buildIsLoading.asObservable();
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
        .finally(() => (this.planStatusesLoading = false));
      this.agentStatusesLoading = true;
      this.envService
        .initAgentStatuses()
        .catch((err) => {
          this.snackbar.open(err, 'Okay', {
            panelClass: ['bg-danger', 'text-white']
          });
        })
        .finally(() => (this.agentStatusesLoading = false));
      this.envService.startAgentStatusSubscription();
      this.envService.startStatusSubscription();
    });
    this.unsubscribe.push(sub1);
    const sub2 = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkBuildStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub2);
    const sub3 = this.envService.agentStatusUpdate.asObservable().subscribe(() => {
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub3);
  }

  ngOnDestroy(): void {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
    this.envService.stopAgentStatusSubscription();
    this.envService.stopStatusSubscription();
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

  rebuildEnv(): void {
    this.isRebuildLoading = true;
    this.cdRef.detectChanges();
    console.log('rebuilding env...');
    this.rebuild
      .executeRebuild()
      .then(
        (success) => {
          if (success) {
            this.isRebuildLoading = false;
            this.router.navigate(['plan']);
          } else {
            this.rebuildErrors = [Error('Rebuild was unsuccessfull, please check server logs for failure point.')];
          }
        },
        (errs) => {
          this.rebuildErrors = errs;
        }
      )
      .finally(() => this.cdRef.detectChanges());
    console.log('done rebuilding env...');
  }

  toggleSelectionMode(): void {
    this.selectionMode = !this.selectionMode;
  }

  toggleDeleteBuildModal(): void {
    this.dialog.open(DeleteBuildModalComponent, {
      width: '50%',
      data: {
        buildName: `${this.envService.getEnvironmentInfo().getValue().name} v${this.envService.getBuildTree().getValue().revision}`,
        buildId: this.envService.getBuildTree().getValue().id
      }
    });
  }

  canDeleteBuild(): boolean {
    return (
      this.buildStatus &&
      (this.buildStatus.state === LaForgeProvisionStatus.Complete ||
        this.buildStatus.state === LaForgeProvisionStatus.Failed ||
        this.buildStatus.state === LaForgeProvisionStatus.Tainted)
    );
  }

  canRebuildBuild(): boolean {
    return (
      this.buildStatus &&
      this.buildStatus.state !== LaForgeProvisionStatus.Planning &&
      this.buildStatus.state !== LaForgeProvisionStatus.Deleted &&
      this.buildStatus.state !== LaForgeProvisionStatus.Todelete &&
      this.buildStatus.state !== LaForgeProvisionStatus.Deleteinprogress &&
      this.buildStatus.state !== LaForgeProvisionStatus.Inprogress
    );
  }
}
