import { Injectable } from '@angular/core';
import {
  LaForgeGetBuildTreeGQL,
  LaForgeGetBuildTreeQuery,
  LaForgeGetEnvironmentGQL,
  LaForgeGetEnvironmentInfoQuery,
  LaForgeGetEnvironmentsQuery,
  LaForgeSubscribeUpdatedAgentStatusGQL,
  LaForgeSubscribeUpdatedAgentStatusSubscription,
  LaForgeSubscribeUpdatedStatusGQL,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { ApiService } from '@services/api/api.service';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class EnvironmentService {
  // private currEnvironment: BehaviorSubject<Environment> = new BehaviorSubject(null);
  private environments: BehaviorSubject<LaForgeGetEnvironmentsQuery['environments']>;
  public envIsLoading: BehaviorSubject<boolean> = new BehaviorSubject(false);
  // private currBuild: BehaviorSubject<Build> = new BehaviorSubject(null);
  // private builds: BehaviorSubject<BuildInfo[]>
  // private agentStatusQuery: QueryRef<AgentStatusQueryResult, EmptyObject>;
  // private agentStatusSubscription: Subscription;
  // private watchingAgentStatus = false;
  // public pollingInterval = 60;
  private environmentInfo: BehaviorSubject<LaForgeGetEnvironmentInfoQuery['environment']>;
  private buildTree: BehaviorSubject<LaForgeGetBuildTreeQuery['build']>;
  private statusMap: { [key: string]: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'] };
  public statusUpdate: BehaviorSubject<boolean>;
  private agentStatusMap: { [key: string]: LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'] };
  public agentStatusUpdate: BehaviorSubject<boolean>;

  constructor(
    private api: ApiService,
    private getEnvironmentInfoGQL: LaForgeGetEnvironmentGQL,
    private getBuildTreeGQL: LaForgeGetBuildTreeGQL,
    private subscribeUpdatedStatus: LaForgeSubscribeUpdatedStatusGQL,
    private subscribeUpdatedAgentStatus: LaForgeSubscribeUpdatedAgentStatusGQL
  ) {
    this.environments = new BehaviorSubject([]);
    this.envIsLoading = new BehaviorSubject(false);
    this.environmentInfo = new BehaviorSubject(null);
    this.buildTree = new BehaviorSubject(null);
    this.statusUpdate = new BehaviorSubject(false);
    this.agentStatusUpdate = new BehaviorSubject(false);
    this.statusMap = {};
    this.agentStatusMap = {};

    this.initEnvironments();
    this.startStatusSubscription();
    this.startAgentStatusSubscription();
    // this.environmentInfo = environmentInfoApolloQuery.watch()
  }

  // public getCurrentEnv(): BehaviorSubject<Environment> {
  //   return this.currEnvironment;
  // }

  public getEnvironmentInfo(): BehaviorSubject<LaForgeGetEnvironmentInfoQuery['environment']> {
    return this.environmentInfo;
  }

  // public getCurrentBuild(): BehaviorSubject<Build> {
  //   return this.currBuild;
  // }

  public getBuildTree(): BehaviorSubject<LaForgeGetBuildTreeQuery['build']> {
    return this.buildTree;
  }

  public getEnvironments(): BehaviorSubject<LaForgeGetEnvironmentsQuery['environments']> {
    return this.environments;
  }

  public getStatus(statusId: string) {
    return this.statusMap[statusId];
  }

  public getAgentStatus(hostId: string) {
    return this.agentStatusMap[hostId];
  }

  private initEnvironments() {
    this.api.pullEnvironments().then((envs) => {
      this.environments.next(envs);
      if (localStorage.getItem('selected_env') && localStorage.getItem('selected_build')) {
        this.setCurrentEnv(localStorage.getItem('selected_env'), localStorage.getItem('selected_build'));
      }
    });
  }

  public initPlanStatuses() {
    if (!this.buildTree.getValue()) return;
    this.api.pullAllPlanStatuses(this.buildTree.getValue().id).then((plans) => {
      for (const plan of plans) {
        this.statusMap[plan.PlanToStatus.id] = { ...plan.PlanToStatus };
      }
      this.statusUpdate.next(!this.statusUpdate.getValue());
    }, console.error);
  }

  public initAgentStatuses() {
    if (!this.buildTree.getValue()) return;
    this.api.pullAllAgentStatuses(this.buildTree.getValue().id).then((build) => {
      for (const team of build.buildToTeam) {
        for (const pnet of team.TeamToProvisionedNetwork) {
          for (const phost of pnet.ProvisionedNetworkToProvisionedHost) {
            if (phost.ProvisionedHostToAgentStatus)
              this.agentStatusMap[phost.ProvisionedHostToAgentStatus.clientId] = {
                ...phost.ProvisionedHostToAgentStatus
              };
          }
        }
      }
      this.agentStatusUpdate.next(!this.agentStatusUpdate.getValue());
    }, console.error);
  }

  public setCurrentEnv(envId: string, buildId: string): void {
    localStorage.setItem('selected_env', `${envId}`);
    localStorage.setItem('selected_build', `${buildId}`);
    // this.pullEnvironment(envId);
    this.pullEnvironmentInfo(envId);
    this.pullBuildTree(buildId);
  }

  public pullEnvironmentInfo(envId: string) {
    this.api.pullEnvironmentInfo(envId).then(
      (env) => {
        if (env?.id) {
          return this.environmentInfo.next(env);
        }
        this.environmentInfo.error(Error('Unable to retrieve environment info. Unknown error.'));
      },
      (err) => {
        localStorage.setItem('selected_env', '');
        this.environmentInfo.error(err);
      }
    );
  }

  public pullBuildTree(buildId: string) {
    this.api.pullBuildTree(buildId).then(
      (build) => {
        if (build?.id) {
          return this.buildTree.next(build);
        }
        this.buildTree.error(Error('Unable to retrieve build tree. Unknown error.'));
      },
      (err) => {
        localStorage.setItem('selected_build', '');
        this.buildTree.error(err);
      }
    );
  }

  private startStatusSubscription() {
    this.subscribeUpdatedStatus.subscribe().subscribe(({ data: { updatedStatus }, errors }) => {
      // console.log('status subscribe');
      if (errors) {
        console.error(errors);
      } else if (updatedStatus) {
        this.statusMap[updatedStatus.id] = {
          ...updatedStatus
        };
        this.statusUpdate.next(!this.statusUpdate.getValue());
      }
    });
  }

  private startAgentStatusSubscription() {
    this.subscribeUpdatedAgentStatus.subscribe().subscribe(({ data: { updatedAgentStatus }, errors }) => {
      if (errors) {
        console.error(errors);
      } else if (updatedAgentStatus) {
        this.agentStatusMap[updatedAgentStatus.clientId] = {
          ...updatedAgentStatus
        };
        this.agentStatusUpdate.next(!this.agentStatusUpdate.getValue());
      }
    });
  }

  // public isWatchingAgentStatus(): boolean {
  //   return this.watchingAgentStatus;
  // }

  // public updateAgentStatuses() {
  //   if (this.currBuild.getValue() === null) return;
  //   const oldBuild = { ...this.currBuild.getValue() };
  //   this.api.pullAgentStatuses(oldBuild.id).then(
  //     (res) => {
  //       console.log(res);
  //       this.currBuild.next(resolveBuildEnums(updateBuildAgentStatuses(oldBuild, res)));
  //     },
  //     () => {
  //       this.currBuild.error(Error('error while pulling agent statuses'));
  //     }
  //   );
  // }

  // public updatePlanStatuses(): void {
  //   if (this.currBuild.getValue() === null) return;
  //   const oldBuild = { ...this.currBuild.getValue() };
  //   this.api.pullBuildPlans(oldBuild.id).then(
  //     (res) => {
  //       this.currBuild.next(updateBuildPlans(oldBuild, res));
  //     },
  //     () => {
  //       this.currBuild.error(Error('error while pulling build plans'));
  //     }
  //   );
  // }

  // public watchAgentStatuses(): void {
  //   this.watchingAgentStatus = true;
  //   if (environment.isMockApi) {
  //     // this.api.pullAgentStatuses(this.currBuild.getValue().id).then(
  //     //   (res) => {
  //     //     this.currEnvironment.next(resolveEnvEnums(updateEnvAgentStatuses(this.currEnvironment.getValue(), res)));
  //     //     this.currBuild.next(resolveBuildEnums(updateBuildAgentStatuses(this.currBuild.getValue(), res)));
  //     //   },
  //     //   (err) => {
  //     //     /* eslint-disable-next-line quotes */
  //     //     this.currEnvironment.error({ ...err, message: "Couldn't load mock data" });
  //     //     this.currBuild.error({ ...err, message: "Couldn't load mock data" });
  //     //     // this.cdRef.detectChanges();
  //     //   }
  //     // );
  //   } else {
  //     this.agentStatusQuery = this.api.getAgentStatuses(this.currBuild.getValue().id);
  //     this.agentStatusQuery.startPolling(this.pollingInterval * 1000);
  //     this.api.setStatusPollingInterval(this.pollingInterval);
  //     // Force UI to refresh so we can detect stale agent data
  //     // this.agentPollingInterval = setInterval(() => this.cdRef.detectChanges(), this.pollingInterval);
  //     this.agentStatusSubscription = this.agentStatusQuery.valueChanges.subscribe(
  //       (res) => {
  //         if (res.data) {
  //           // const updatedEnv = resolveEnvEnums(updateEnvAgentStatuses(this.currEnvironment.getValue(), res.data.));
  //           const updatedBuild = resolveBuildEnums(updateBuildAgentStatuses(this.currBuild.getValue(), res.data));
  //           // this.currEnvironment.next(updatedEnv);
  //           this.currBuild.next(updatedBuild);
  //         }
  //       },
  //       (err) => {
  //         // this.currEnvironment.error({ ...err, message: 'Too many database connections' });
  //         this.currBuild.error({ ...err, message: 'Too many database connections' });
  //       }
  //     );
  //   }
  // }

  // public setAgentPollingInterval(interval: number) {
  //   this.pollingInterval = interval;
  //   if (this.agentStatusQuery) {
  //     this.agentStatusQuery.stopPolling();
  //     this.agentStatusQuery.startPolling(interval * 1000);
  //   }
  //   this.api.setStatusPollingInterval(interval);
  //   // this.agentPollingInterval = setInterval(() => this.cdRef.detectChanges(), this.pollingInterval);
  //   // this.cdRef.detectChanges();
  // }

  // public stopWatchingAgentStatus(): void {
  //   this.agentStatusQuery.stopPolling();
  //   this.agentStatusQuery = null;
  //   this.watchingAgentStatus = false;
  //   this.agentStatusSubscription.unsubscribe();
  // }

  // private pullEnvironment(id: ID) {
  //   this.envIsLoading.next(true);
  //   this.api.pullEnvironment(id).then(
  //     (env: Environment) => {
  //       // env = resolveEnvEnums(env);
  //       // this.currEnvironment.next({
  //       //   ...env,
  //       //   EnvironmentToBuild: [...env.EnvironmentToBuild]
  //       //     .sort((a, b) => b.revision - a.revision)
  //       //     .map((build) => ({
  //       //       ...build,
  //       //       buildToTeam: [...build.buildToTeam]
  //       //         .sort((a, b) => a.team_number - b.team_number)
  //       //         .map((team) => ({
  //       //           ...team,
  //       //           TeamToProvisionedNetwork: [...team.TeamToProvisionedNetwork]
  //       //             .sort((a, b) => {
  //       //               if (a.name < b.name) return -1;
  //       //               if (a.name > b.name) return 1;
  //       //               return 0;
  //       //             })
  //       //             .map((network) => ({
  //       //               ...network,
  //       //               ProvisionedNetworkToProvisionedHost: [...network.ProvisionedNetworkToProvisionedHost].sort((a, b) => {
  //       //                 if (a.ProvisionedHostToHost.hostname < b.ProvisionedHostToHost.hostname) return -1;
  //       //                 if (a.ProvisionedHostToHost.hostname > b.ProvisionedHostToHost.hostname) return 1;
  //       //                 return 0;
  //       //               })
  //       //             }))
  //       //         }))
  //       //     }))
  //       // });
  //       this.currEnvironment.next(env);
  //       if (localStorage.getItem('selected_build')) {
  //         // const builds = env.EnvironmentToBuild.filter((build) => build.id === localStorage.getItem('selected_build'));
  //         // if (builds.length === 1) {
  //         //   this.currBuild.next(
  //         //     resolveBuildEnums({
  //         //       ...builds[0],
  //         //       buildToTeam: [...builds[0].buildToTeam]
  //         //         .sort((a, b) => a.team_number - b.team_number)
  //         //         .map((team) => ({
  //         //           ...team,
  //         //           TeamToProvisionedNetwork: [...team.TeamToProvisionedNetwork]
  //         //             .sort((a, b) => {
  //         //               if (a.name < b.name) return -1;
  //         //               if (a.name > b.name) return 1;
  //         //               return 0;
  //         //             })
  //         //             .map((network) => ({
  //         //               ...network,
  //         //               ProvisionedNetworkToProvisionedHost: [...network.ProvisionedNetworkToProvisionedHost].sort((a, b) => {
  //         //                 if (a.ProvisionedHostToHost.hostname < b.ProvisionedHostToHost.hostname) return -1;
  //         //                 if (a.ProvisionedHostToHost.hostname > b.ProvisionedHostToHost.hostname) return 1;
  //         //                 return 0;
  //         //               })
  //         //             }))
  //         //         }))
  //         //     })
  //         //   );
  //         // } else {
  // eslint-disable-next-line max-len
  //         //   this.currBuild.error(Error(`error locating the selected build, found ${builds.length} builds that match the selected id`));
  //         // }
  //         this.pullBuild(localStorage.getItem('selected_build'));
  //       } else {
  //         this.envIsLoading.next(false);
  //       }
  //     },
  //     (err) => {
  //       this.currEnvironment.error(err);
  //     }
  //   );
  // }

  // private pullBuild(id: ID) {
  //   this.api
  //     .pullBuild(id)
  //     .then(
  //       (build) => {
  //         this.currBuild.next(
  //           resolveBuildEnums({
  //             ...build,
  //             buildToTeam: [...build.buildToTeam]
  //               .sort((a, b) => a.team_number - b.team_number)
  //               .map((team) => ({
  //                 ...team,
  //                 TeamToProvisionedNetwork: [...team.TeamToProvisionedNetwork]
  //                   .sort((a, b) => {
  //                     if (a.name < b.name) return -1;
  //                     if (a.name > b.name) return 1;
  //                     return 0;
  //                   })
  //                   .map((network) => ({
  //                     ...network,
  //                     ProvisionedNetworkToProvisionedHost: [...network.ProvisionedNetworkToProvisionedHost]
  //                       .sort((a, b) => {
  //                         if (a.ProvisionedHostToHost.hostname < b.ProvisionedHostToHost.hostname) return -1;
  //                         if (a.ProvisionedHostToHost.hostname > b.ProvisionedHostToHost.hostname) return 1;
  //                         return 0;
  //                       })
  //                       .map((host) => ({
  //                         ...host,
  //                         ProvisionedHostToProvisioningStep: [...host.ProvisionedHostToProvisioningStep].sort((a, b) => {
  //                           if (a.step_number !== null && b.step_number !== null) return a.step_number - b.step_number;
  //                           return 0;
  //                         })
  //                       }))
  //                   }))
  //               }))
  //           })
  //         );
  //       },
  //       (err) => {
  //         this.currBuild.error(Error('error ocurred while loading build'));
  //       }
  //     )
  //     .finally(() => {
  //       this.envIsLoading.next(false);
  //     });
  // }
}
