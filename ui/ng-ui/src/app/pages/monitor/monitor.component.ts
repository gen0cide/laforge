import { ChangeDetectorRef, Component, OnDestroy, AfterViewInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { LaForgeGetBuildTreeQuery, LaForgeGetEnvironmentInfoQuery } from '@graphql';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { Observable, Subscription } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { AgentStatusQueryResult } from 'src/app/models/api.model';
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
  environment: Observable<LaForgeGetEnvironmentInfoQuery['environment']> = null;
  envLoaded = false;
  build: Observable<LaForgeGetBuildTreeQuery['build']>;
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
    this.subheader.setShowEnvDropdown(true);

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
  }

  ngAfterViewInit(): void {
    // this.environment.subscribe((environment) => {
    //   if (environment && !this.envService.isWatchingAgentStatus()) {
    //     setTimeout(() => this.envService.watchAgentStatuses(), 1000);
    //   }
    // });
    // this.build.subscribe((build) => {
    //   if (build && !this.envService.isWatchingAgentStatus()) {
    //     setTimeout(() => this.envService.watchAgentStatuses(), 1000);
    //   }
    // });
  }

  envIsSelected(): boolean {
    return this.envService.getEnvironmentInfo() != null && environment != null;
  }

  ngOnDestroy(): void {}

  onIntervalChange(changeEvent: MatSelectChange): void {
    // this.envService.setAgentPollingInterval(changeEvent.value);
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }
}
