import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { Subscription } from 'rxjs';
import { updateAgentStatuses } from 'src/app/models/agent.model';
import { AgentStatusQueryResult, EnvironmentInfo } from 'src/app/models/api.model';
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
export class ManageComponent implements OnInit {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  envs: EnvironmentInfo[];
  environment: Environment = null;
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
  }

  ngOnInit(): void {
    this.api.pullEnvironments().then((envs: EnvironmentInfo[]) => {
      this.envs = envs;
      this.cdRef.detectChanges();
    });
  }

  envIsSelected(): boolean {
    return this.envService.getCurrentEnv() != null;
  }

  grabEnvironmentTree(changeEvent: MatSelectChange): void {
    this.envService.setCurrentEnv(this.envs.filter((e) => e.id === changeEvent.value)[0]);
    this.api.pullEnvTree(this.envService.getCurrentEnv().id).then(
      (env: Environment) => {
        this.environment = {
          ...env,
          build: {
            ...env.build,
            teams: [...env.build.teams]
              .sort((a, b) => a.teamNumber - b.teamNumber)
              .map((team) => ({
                ...team,
                provisionedNetworks: [...team.provisionedNetworks]
                  .sort((a, b) => {
                    if (a.name < b.name) return -1;
                    if (a.name > b.name) return 1;
                    return 0;
                  })
                  .map((network) => ({
                    ...network,
                    provisionedHosts: [...network.provisionedHosts].sort((a, b) => {
                      if (a.host.hostname < b.host.hostname) return -1;
                      if (a.host.hostname > b.host.hostname) return 1;
                      return 0;
                    })
                  }))
              }))
          }
        };
        this.envLoaded = true;
        this.cdRef.detectChanges();
        this.initAgentStatusPolling();
      },
      (err) => {
        this.apolloError = err;
        this.cdRef.detectChanges();
        // console.log(typeof err);
        // console.log(err.toString());
        // console.error('yep, cant connect');
        // console.error(err);
      }
    );
  }

  initAgentStatusPolling(): void {
    if (environment.isMockApi) {
      this.api.pullAgentStatuses(this.environment.id).then(
        (res) => {
          this.environment = updateAgentStatuses(this.environment, res);
          this.loading = false;
          this.apolloError = {};
          this.cdRef.detectChanges();
        },
        (err) => {
          /* eslint-disable-next-line quotes */
          this.apolloError = { ...err, message: "Couldn't load mock data" };
          this.cdRef.detectChanges();
        }
      );
      return;
    }
    console.log('Agent status polling initializing...');
    this.agentStatusQuery = this.api.getAgentStatuses(this.environment.id);
    this.agentStatusQuery.startPolling(this.pollingInterval * 1000);
    this.api.setStatusPollingInterval(this.pollingInterval);
    // Force UI to refresh so we can detect stale agent data
    this.agentPollingInterval = setInterval(() => this.cdRef.detectChanges(), this.pollingInterval);
    this.agentStatusSubscription = this.agentStatusQuery.valueChanges.subscribe(
      ({ data: result }) => {
        if (result) {
          this.loading = false;
          this.environment = updateAgentStatuses(this.environment, result);
          this.apolloError = {};
          // console.log('data updated');
        }
      },
      (err) => {
        this.apolloError = { ...err, message: 'Too many database connections' };
        this.cdRef.detectChanges();
      }
    );
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }

  toggleSelectionMode(): void {
    this.selectionMode = !this.selectionMode;
  }
}
