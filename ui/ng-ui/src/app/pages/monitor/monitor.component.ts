import { ChangeDetectorRef, Component, OnInit, OnDestroy } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { Subscription } from 'rxjs';
import { updateAgentStatuses } from 'src/app/models/agent.model';
import { AgentStatusQueryResult, EnvironmentInfo } from 'src/app/models/api.model';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { filter } from 'rxjs/operators';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { ApolloError } from '@apollo/client/core';

@Component({
  selector: 'app-manage',
  templateUrl: './monitor.component.html',
  styleUrls: ['./monitor.component.scss']
})
export class MonitorComponent implements OnInit, OnDestroy {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  envs: EnvironmentInfo[];
  environment: Environment = null;
  envLoaded = false;
  environmentDetailsCols: string[] = ['TeamCount', 'AdminCIDRs', 'ExposedVDIPorts', 'maintainer'];
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
    private envService: EnvironmentService
  ) {
    this.subheader.setTitle('Monitor Agents');
    this.subheader.setDescription('View live data being sent from the host agents');
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
            teams: [...env.build.teams].sort((a,b) => a.teamNumber - b.teamNumber).map(team => ({
              ...team,
              provisionedNetworks: [...team.provisionedNetworks].sort((a, b) => {
                if (a.name < b.name) return -1;
                if (a.name > b.name) return 1;
                return 0;
              }).map(network => ({
                ...network,
                provisionedHosts: [...network.provisionedHosts].sort((a, b) => {
                  if (a.host.hostname < b.host.hostname) return -1;
                  if (a.host.hostname > b.host.hostname) return 1;
                  return 0;
                })
              }))
            }))
          }
        }
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

  ngOnDestroy(): void {
    clearInterval(this.agentPollingInterval);
    this.agentStatusSubscription.unsubscribe();
  }

  initAgentStatusPolling(): void {
    console.log('Agent status polling initializing...');
    this.agentStatusQuery = this.api.getAgentStatuses(this.environment.id);
    this.agentStatusQuery.startPolling(this.pollingInterval * 1000);
    this.api.setStatusPollingInterval(this.pollingInterval);
    // Force UI to refresh so we can detect stale agent data
    this.agentPollingInterval = setInterval(() => this.cdRef.detectChanges(), this.pollingInterval);
    this.agentStatusSubscription = this.agentStatusQuery.valueChanges.subscribe(({ data: result }) => {
      if (result) {
        this.loading = false;
        this.environment = updateAgentStatuses(this.environment, result);
        this.apolloError = {};
        // console.log('data updated');
      }
    }, (err) => {
      this.apolloError = { ...err, message: 'Too many database connections' };
      this.cdRef.detectChanges();
    });
  }

  fetchAgentStatuses(): void {
    // console.log('Polling agent statuses...');
    this.loading = true;
    this.cdRef.detectChanges();
    // this.api.getAgentStatuses(this.environment.id).then((result: AgentStatusQueryResult) => {
    //   this.loading = false;
    //   this.environment = updateAgentStatuses(this.environment, result);
    //   this.cdRef.detectChanges();
    // }, console.error);
  }
  // onBranchSelect(changeEvent: MatSelectChange) {
  onIntervalChange(changeEvent: MatSelectChange): void {
    // Update the interval based on select's value
    this.pollingInterval = changeEvent.value;
    this.agentStatusQuery.stopPolling();
    this.agentStatusQuery.startPolling(changeEvent.value * 1000);
    this.api.setStatusPollingInterval(this.pollingInterval);
    this.agentPollingInterval = setInterval(() => this.cdRef.detectChanges(), this.pollingInterval);
    this.cdRef.detectChanges();
    // Stop the old polling
    // clearInterval(this.agentPollingInterval);
    // Set up polling again with new interval
    // this.agentPollingInterval = setInterval(() => this.fetchAgentStatuses(), this.pollingInterval * 1000);
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }
}
