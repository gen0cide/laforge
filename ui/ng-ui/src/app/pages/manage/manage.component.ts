import { AfterViewInit, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { Observable, Subscription } from 'rxjs';
import { updateAgentStatuses } from 'src/app/models/agent.model';
import { AgentStatusQueryResult, EnvironmentInfo } from 'src/app/models/api.model';
import { ID } from 'src/app/models/common.model';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';
import { environment } from 'src/environments/environment';
import { SubheaderService } from '../../_metronic/partials/layout/subheader/_services/subheader.service';

@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ManageComponent implements AfterViewInit {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  envs: EnvironmentInfo[];
  environment: Observable<Environment>;
  envIsLoading: Observable<boolean>;
  loaded = false;
  environmentDetailsCols: string[] = ['TeamCount', 'AdminCIDRs', 'ExposedVDIPorts', 'maintainer'];
  agentPollingInterval: NodeJS.Timeout;
  pollingInterval = 60;
  loading = false;
  intervalOptions = [10, 30, 60, 120];
  selectionMode = false;
  envLoaded = false;
  agentStatusQuery: QueryRef<AgentStatusQueryResult, EmptyObject>;
  agentStatusSubscription: Subscription;
  apolloError: any = {};

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    private envService: EnvironmentService
  ) {
    this.subheader.setTitle('Environment');
    this.subheader.setDescription('Manage your currently running environment');

    this.environment = this.envService.getCurrentEnv().asObservable();
    this.envIsLoading = this.envService.envIsLoading.asObservable();
  }

  ngAfterViewInit(): void {
    this.environment.subscribe((environment) => {
      if (environment && !this.envService.isWatchingAgentStatus()) {
        setTimeout(() => this.envService.watchAgentStatuses(), 1000);
      }
    });
  }

  envIsSelected(): boolean {
    return this.envService.getCurrentEnv().getValue() != null;
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }

  toggleSelectionMode(): void {
    this.selectionMode = !this.selectionMode;
  }
}
