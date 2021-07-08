import { ChangeDetectorRef, Component, OnInit, OnDestroy, AfterViewInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { ApolloError } from '@apollo/client/core';
import { RebuildService } from '@services/rebuild/rebuild.service';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { Observable, Subscription } from 'rxjs';
import { filter } from 'rxjs/operators';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { updateEnvAgentStatuses } from 'src/app/models/agent.model';
import { AgentStatusQueryResult, EnvironmentInfo } from 'src/app/models/api.model';
import { ID } from 'src/app/models/common.model';
import { Build, Environment, resolveEnvEnums } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-manage',
  templateUrl: './monitor.component.html',
  styleUrls: ['./monitor.component.scss']
})
export class MonitorComponent implements AfterViewInit, OnDestroy {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  // envs: EnvironmentInfo[];
  environment: Observable<Environment> = null;
  envLoaded = false;
  build: Observable<Build>;
  environmentDetailsCols: string[] = ['TeamCount', 'AdminCIDRs', 'ExposedVDIPorts'];
  agentPollingInterval: NodeJS.Timeout;
  pollingInterval = 60;
  loading = false;
  intervalOptions = [10, 30, 60, 120];
  agentStatusQuery: QueryRef<AgentStatusQueryResult, EmptyObject>;
  agentStatusSubscription: Subscription;
  apolloError: any = {};

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    public envService: EnvironmentService
  ) {
    this.subheader.setTitle('Monitor Agents');
    this.subheader.setDescription('View live data being sent from the host agents');

    this.environment = this.envService.getCurrentEnv().asObservable();
    this.build = this.envService.getCurrentBuild().asObservable();
  }

  ngAfterViewInit(): void {
    // this.environment.subscribe((environment) => {
    //   if (environment && !this.envService.isWatchingAgentStatus()) {
    //     setTimeout(() => this.envService.watchAgentStatuses(), 1000);
    //   }
    // });
    this.build.subscribe((build) => {
      if (build && !this.envService.isWatchingAgentStatus()) {
        setTimeout(() => this.envService.watchAgentStatuses(), 1000);
      }
    });
  }

  envIsSelected(): boolean {
    return this.envService.getCurrentEnv() != null && environment != null;
  }

  ngOnDestroy(): void {}

  onIntervalChange(changeEvent: MatSelectChange): void {
    this.envService.setAgentPollingInterval(changeEvent.value);
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }
}
