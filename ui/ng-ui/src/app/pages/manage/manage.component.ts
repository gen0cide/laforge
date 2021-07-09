import { AfterViewInit, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { LaForgeEnvironment, LaForgeGetEnvironmentInfoQuery, LaForgeGetBuildTreeQuery } from '@graphql';
import { RebuildService } from '@services/rebuild/rebuild.service';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { GraphQLError } from 'graphql';
import { interval, Observable, Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { updateEnvAgentStatuses } from 'src/app/models/agent.model';
import { AgentStatusQueryResult, EnvironmentInfo } from 'src/app/models/api.model';
import { ID } from 'src/app/models/common.model';
import { Build, Environment, resolveEnvEnums } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

import { SubheaderService } from '../../_metronic/partials/layout/subheader/_services/subheader.service';

@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ManageComponent implements OnInit {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  envs: EnvironmentInfo[];
  environment: Observable<LaForgeGetEnvironmentInfoQuery['environment']>;
  build: Observable<LaForgeGetBuildTreeQuery['build']>;
  envIsLoading: Observable<boolean>;
  loaded = false;
  environmentDetailsCols: string[] = ['TeamCount', 'AdminCIDRs', 'ExposedVDIPorts'];
  agentPollingInterval: NodeJS.Timeout;
  pollingInterval = 60;
  loading = false;
  intervalOptions = [10, 30, 60, 120];
  selectionMode = false;
  envLoaded = false;
  agentStatusQuery: QueryRef<AgentStatusQueryResult, EmptyObject>;
  agentStatusSubscription: Subscription;
  apolloError: any = {};
  isRebuildLoading = false;
  rebuildErrors: (GraphQLError | Error)[] = [];

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    private envService: EnvironmentService,
    private rebuild: RebuildService
  ) {
    this.subheader.setTitle('Environment');
    this.subheader.setDescription('Manage your currently running environment');

    // this.environment = this.envService.getCurrentEnv().asObservable();
    // this.build = this.envService.getCurrentBuild().asObservable();
    // this.envIsLoading = this.envService.envIsLoading.asObservable();
    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
  }

  ngOnInit(): void {
    // pull the statuses on load and then poll every 10 secs
    // interval(10000).subscribe(() => {
    //   this.envService.updateAgentStatuses();
    // });
  }

  envIsSelected(): boolean {
    return this.envService.getEnvironmentInfo().getValue() != null;
  }

  rebuildEnv(): void {
    //   this.isRebuildLoading = true;
    //   this.cdRef.detectChanges();
    //   console.log('rebuilding env...');
    //   this.rebuild
    //     .executeRebuild()
    //     .then(
    //       (success) => {
    //         if (success) {
    //           this.isRebuildLoading = false;
    //         } else {
    //           this.rebuildErrors = [Error('Rebuild was unsuccessfull, please check server logs for failure point.')];
    //         }
    //       },
    //       (errs) => {
    //         this.rebuildErrors = errs;
    //       }
    //     )
    //     .finally(() => this.cdRef.detectChanges());
    //   console.log('done rebuilding env...');
  }

  toggleSelectionMode(): void {
    //   this.selectionMode = !this.selectionMode;
  }
}
