import { Injectable } from '@angular/core';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { updateAgentStatuses } from 'src/app/models/agent.model';
import { AgentStatusQueryResult, EnvironmentInfo } from 'src/app/models/api.model';
import { ID } from 'src/app/models/common.model';
import { Environment } from 'src/app/models/environment.model';
import { environment } from 'src/environments/environment';
import { ApiService } from '../api/api.service';

@Injectable({
  providedIn: 'root'
})
export class EnvironmentService {
  private currEnvironment: BehaviorSubject<Environment> = new BehaviorSubject(null);
  private environments: BehaviorSubject<EnvironmentInfo[]> = new BehaviorSubject([]);
  public envIsLoading: BehaviorSubject<boolean> = new BehaviorSubject(false);
  private agentStatusQuery: QueryRef<AgentStatusQueryResult, EmptyObject>;
  private agentStatusSubscription: Subscription;
  private watchingAgentStatus = false;
  public pollingInterval = 60;

  constructor(private api: ApiService) {
    this.initEnvironments();
  }

  public getCurrentEnv(): BehaviorSubject<Environment> {
    return this.currEnvironment;
  }

  public setCurrentEnv(id: ID): void {
    localStorage.setItem('selected_env', `${id}`);
    this.pullEnvironment(id);
  }

  public getEnvironments(): BehaviorSubject<EnvironmentInfo[]> {
    return this.environments;
  }

  public isWatchingAgentStatus(): boolean {
    return this.watchingAgentStatus;
  }

  public watchAgentStatuses(): void {
    this.watchingAgentStatus = true;
    if (environment.isMockApi) {
      this.api.pullAgentStatuses(this.currEnvironment.getValue().id).then(
        (res) => {
          this.currEnvironment.next(updateAgentStatuses(this.currEnvironment.getValue(), res));
        },
        (err) => {
          /* eslint-disable-next-line quotes */
          this.currEnvironment.error({ ...err, message: "Couldn't load mock data" });
          // this.cdRef.detectChanges();
        }
      );
    } else {
      this.agentStatusQuery = this.api.getAgentStatuses(this.currEnvironment.getValue().id);
      this.agentStatusQuery.startPolling(this.pollingInterval * 1000);
      this.api.setStatusPollingInterval(this.pollingInterval);
      // Force UI to refresh so we can detect stale agent data
      // this.agentPollingInterval = setInterval(() => this.cdRef.detectChanges(), this.pollingInterval);
      this.agentStatusSubscription = this.agentStatusQuery.valueChanges.subscribe(
        ({ data: result }) => {
          if (result) {
            this.currEnvironment.next(updateAgentStatuses(this.currEnvironment.getValue(), result));
            // console.log('data updated');
          }
        },
        (err) => {
          this.currEnvironment.error({ ...err, message: 'Too many database connections' });
        }
      );
    }
  }

  public setAgentPollingInterval(interval: number) {
    this.pollingInterval = interval;
    if (this.agentStatusQuery) {
      this.agentStatusQuery.stopPolling();
      this.agentStatusQuery.startPolling(interval * 1000);
    }
    this.api.setStatusPollingInterval(interval);
    // this.agentPollingInterval = setInterval(() => this.cdRef.detectChanges(), this.pollingInterval);
    // this.cdRef.detectChanges();
  }

  public stopWatchingAgentStatus(): void {
    this.agentStatusQuery.stopPolling();
    this.agentStatusQuery = null;
    this.watchingAgentStatus = false;
    this.agentStatusSubscription.unsubscribe();
  }

  private initEnvironments() {
    this.api.pullEnvironments().then((envs) => {
      this.environments.next(envs);
      if (localStorage.getItem('selected_env')) {
        this.setCurrentEnv(localStorage.getItem('selected_env'));
      }
    });
  }

  private pullEnvironment(id: ID) {
    this.envIsLoading.next(true);
    this.api.pullEnvTree(id).then(
      (env: Environment) => {
        this.currEnvironment.next({
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
        });
        this.envIsLoading.next(false);
      },
      (err) => {
        this.currEnvironment.error(err);
      }
    );
  }
}
