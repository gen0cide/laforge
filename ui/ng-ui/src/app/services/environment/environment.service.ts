import { Injectable } from '@angular/core';
import {
  LaForgeAgentStatusFieldsFragment,
  LaForgeAgentTaskFieldsFragment,
  LaForgeBuildCommitFieldsFragment,
  LaForgeGetBuildTreeGQL,
  LaForgeGetBuildTreeQuery,
  LaForgeGetEnvironmentGQL,
  LaForgeGetEnvironmentInfoQuery,
  LaForgeGetEnvironmentsQuery,
  LaForgePlanFieldsFragment,
  LaForgeStatusFieldsFragment,
  LaForgeSubscribeUpdatedAgentStatusGQL,
  LaForgeSubscribeUpdatedAgentStatusSubscription,
  LaForgeSubscribeUpdatedAgentTaskGQL,
  LaForgeSubscribeUpdatedBuildCommitGQL,
  LaForgeSubscribeUpdatedBuildGQL,
  LaForgeSubscribeUpdatedStatusGQL,
  LaForgeSubscribeUpdatedStatusSubscription
} from '@graphql';
import { ApiService } from '@services/api/api.service';
import { BehaviorSubject, Subscription } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class EnvironmentService {
  // List of envs
  private environments: BehaviorSubject<LaForgeGetEnvironmentsQuery['environments']>;
  // Currently selected data
  private environmentInfo: BehaviorSubject<LaForgeGetEnvironmentInfoQuery['environment']>;
  private buildTree: BehaviorSubject<LaForgeGetBuildTreeQuery['build']>;
  // Loading
  public envIsLoading: BehaviorSubject<boolean> = new BehaviorSubject(false);
  public buildIsLoading: BehaviorSubject<boolean> = new BehaviorSubject(false);
  // Status
  private statusMap: { [key: string]: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'] };
  public statusUpdate: BehaviorSubject<boolean>;
  private statusSubscription: Subscription;
  // Agent Status
  private agentStatusMap: { [key: string]: LaForgeSubscribeUpdatedAgentStatusSubscription['updatedAgentStatus'] };
  public agentStatusUpdate: BehaviorSubject<boolean>;
  private agentStatusSubscription: Subscription;
  // Plans
  private planMap: { [key: string]: LaForgePlanFieldsFragment };
  public planUpdate: BehaviorSubject<boolean>;
  // Build
  private buildSubscription: Subscription;
  // Build Commits
  private buildCommitMap: { [key: string]: LaForgeBuildCommitFieldsFragment };
  public buildCommitUpdate: BehaviorSubject<boolean>;
  private buildCommitSubscription: Subscription;
  // Agent Tasks
  private agentTaskMap: { [key: string]: LaForgeAgentTaskFieldsFragment };
  public agentTaskUpdate: BehaviorSubject<boolean>;
  private agentTaskSubscription: Subscription;

  constructor(
    private api: ApiService,
    private getEnvironmentInfoGQL: LaForgeGetEnvironmentGQL,
    private getBuildTreeGQL: LaForgeGetBuildTreeGQL,
    private subscribeUpdatedStatus: LaForgeSubscribeUpdatedStatusGQL,
    private subscribeUpdatedAgentStatus: LaForgeSubscribeUpdatedAgentStatusGQL,
    private subscribeUpdatedBuild: LaForgeSubscribeUpdatedBuildGQL,
    private subscribeUpdatedBuildCommit: LaForgeSubscribeUpdatedBuildCommitGQL,
    private subscribeUpdatedAgentTask: LaForgeSubscribeUpdatedAgentTaskGQL
  ) {
    this.environments = new BehaviorSubject([]);

    this.environmentInfo = new BehaviorSubject(null);
    this.buildTree = new BehaviorSubject(null);

    this.envIsLoading = new BehaviorSubject(false);
    this.buildIsLoading = new BehaviorSubject(false);

    this.statusUpdate = new BehaviorSubject(false);
    this.agentStatusUpdate = new BehaviorSubject(false);
    this.planUpdate = new BehaviorSubject(false);
    this.buildCommitUpdate = new BehaviorSubject(false);
    this.agentTaskUpdate = new BehaviorSubject(false);

    this.statusMap = {};
    this.agentStatusMap = {};
    this.planMap = {};
    this.buildCommitMap = {};
    this.agentTaskMap = {};

    this.initEnvironments();
    // this.startStatusSubscription();
    // this.startAgentStatusSubscription();
    // this.startBuildSubscription();
    // this.startBuildCommitSubscription();
    // this.environmentInfo = environmentInfoApolloQuery.watch()
  }

  public getEnvironmentInfo(): BehaviorSubject<LaForgeGetEnvironmentInfoQuery['environment']> {
    return this.environmentInfo;
  }

  public getBuildTree(): BehaviorSubject<LaForgeGetBuildTreeQuery['build']> {
    return this.buildTree;
  }

  public getEnvironments(): BehaviorSubject<LaForgeGetEnvironmentsQuery['environments']> {
    return this.environments;
  }

  public getStatus(statusId: string): LaForgeStatusFieldsFragment {
    return this.statusMap[statusId];
  }

  public getAgentStatus(hostId: string): LaForgeAgentStatusFieldsFragment {
    return this.agentStatusMap[hostId];
  }

  public getPlan(planId: string): LaForgePlanFieldsFragment {
    return this.planMap[planId];
  }

  public getBuildCommit(buildCommitId: string): LaForgeBuildCommitFieldsFragment {
    return this.buildCommitMap[buildCommitId];
  }

  public getLatestCommit(): LaForgeBuildCommitFieldsFragment | null {
    if (!this.buildTree.getValue()) return null;
    return this.buildCommitMap[this.buildTree.getValue().BuildToLatestBuildCommit.id];
  }

  public getAgentTask(agentTaskId: string): LaForgeAgentTaskFieldsFragment {
    return this.agentTaskMap[agentTaskId];
  }

  public initEnvironments() {
    this.api.pullEnvironments().then((envs) => {
      this.environments.next(envs);
      if (localStorage.getItem('selected_env') && localStorage.getItem('selected_build')) {
        console.log(`currently selected build: ${localStorage.getItem('selected_build')}`);
        this.setCurrentEnv(localStorage.getItem('selected_env'), localStorage.getItem('selected_build'));
      }
    });
  }

  public initPlanStatuses(): Promise<boolean> {
    return new Promise(async (resolve, reject) => {
      if (!this.buildTree.getValue()) return reject(new Error("Can't load Plan Statuses as Build Tree hasn't been loaded."));
      const count = 100;
      let offset = 0;
      let total = 1;
      while (offset < total) {
        console.log(`Getting status | count: ${count}, offset: ${offset}`);
        await this.api.pullAllPlanStatuses(this.buildTree.getValue().id, count, offset).then(
          (data) => {
            for (const status of data.statuses) {
              this.statusMap[status.id] = { ...status };
            }
            this.statusUpdate.next(!this.statusUpdate.getValue());
            total = data.pageInfo.total;
            offset = data.pageInfo.nextOffset;
          },
          (err) => {
            console.error(err);
            reject(err);
          }
        );
      }
      resolve(true);
    });
  }

  public initAgentStatuses(): Promise<boolean> {
    return new Promise(async (resolve, reject) => {
      if (!this.buildTree.getValue()) return reject(new Error("Can't load Agent Statuses as Build Tree hasn't been loaded."));
      const count = 50;
      let offset = 0;
      let total = 1;
      while (offset < total) {
        console.log(`Getting agent status | count: ${count}, offset: ${offset}`);
        await this.api.pullAllAgentStatuses(this.buildTree.getValue().id, count, offset).then(
          (data) => {
            let updatedAgentStatuses = 0;
            for (const agentStatus of data.agentStatuses) {
              if (!this.agentStatusMap[agentStatus.clientId]) {
                this.agentStatusMap[agentStatus.clientId] = {
                  ...agentStatus
                };
                updatedAgentStatuses++;
              }
            }
            if (updatedAgentStatuses) this.agentStatusUpdate.next(!this.agentStatusUpdate.getValue());
            total = data.pageInfo.total;
            offset = data.pageInfo.nextOffset;
          },
          (err) => {
            console.error(err);
            reject(err);
          }
        );
      }
      resolve(true);
    });
  }

  public initPlans(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      if (!this.buildTree.getValue()) return reject(new Error("Can't load Plans as Build Tree hasn't been loaded."));
      this.api.pullBuildPlans(this.buildTree.getValue().id).then(
        (build) => {
          for (const plan of build.buildToPlan) {
            this.planMap[plan.id] = {
              ...plan
            };
          }
          this.planUpdate.next(!this.planUpdate.getValue());
          resolve(true);
        },
        (err) => {
          console.error(err);
          reject(err);
        }
      );
    });
  }

  public initBuildCommits() {
    if (!this.buildTree.getValue()) return;
    this.api.pullBuildCommits(this.buildTree.getValue().id).then((buildCommits) => {
      for (const buildCommit of buildCommits) {
        this.buildCommitMap[buildCommit.id] = {
          ...buildCommit
        };
      }
      this.buildCommitUpdate.next(!this.buildCommitUpdate.getValue());
    });
  }

  public setCurrentEnv(envId: string, buildId: string): void {
    localStorage.setItem('selected_env', `${envId}`);
    localStorage.setItem('selected_build', `${buildId}`);
    this.pullEnvironmentInfo(envId);
    this.pullBuildTree(buildId);
    console.log(`currently selected build: ${localStorage.getItem('selected_build')}`);
  }

  public pullEnvironmentInfo(envId: string) {
    this.envIsLoading.next(true);
    this.api
      .pullEnvironmentInfo(envId)
      .then(
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
      )
      .finally(() => this.envIsLoading.next(false));
  }

  public pullBuildTree(buildId: string) {
    this.buildIsLoading.next(true);
    this.api
      .pullBuildTree(buildId)
      .then(
        (build) => {
          if (build?.id) {
            this.statusMap[build.buildToStatus.id] = { ...build.buildToStatus };
            this.statusUpdate.next(!this.statusUpdate.getValue());
            return this.buildTree.next(build);
          }
          this.buildTree.error(Error('Unable to retrieve build tree. Unknown error.'));
        },
        (err) => {
          localStorage.setItem('selected_build', '');
          this.buildTree.error(err);
        }
      )
      .finally(() => this.buildIsLoading.next(false));
  }

  public startStatusSubscription() {
    this.statusSubscription = this.subscribeUpdatedStatus.subscribe().subscribe(({ data: { updatedStatus }, errors }) => {
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

  public stopStatusSubscription(): void {
    if (this.statusSubscription) this.statusSubscription.unsubscribe();
  }

  public startAgentStatusSubscription() {
    this.agentStatusSubscription = this.subscribeUpdatedAgentStatus.subscribe().subscribe(({ data: { updatedAgentStatus }, errors }) => {
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

  public stopAgentStatusSubscription(): void {
    if (this.agentStatusSubscription) this.agentStatusSubscription.unsubscribe();
  }

  public startBuildSubscription() {
    this.buildSubscription = this.subscribeUpdatedBuild.subscribe().subscribe(({ data: { updatedBuild }, errors }) => {
      if (errors) {
        console.error(errors);
      } else if (updatedBuild) {
        // const oldBuildTree = this.buildTree.getValue();
        this.pullBuildTree(updatedBuild.id);
        // this.buildTree.next({
        //   ...oldBuildTree,
        //   BuildToLatestBuildCommit: {
        //     ...updatedBuild.BuildToLatestBuildCommit
        //   }
        // });
      }
    });
  }

  public stopBuildSubscription(): void {
    if (this.buildSubscription) this.buildSubscription.unsubscribe();
  }

  public startBuildCommitSubscription() {
    this.buildCommitSubscription = this.subscribeUpdatedBuildCommit.subscribe().subscribe(({ data: { updatedCommit }, errors }) => {
      if (errors) {
        console.error(errors);
      } else if (updatedCommit) {
        this.buildCommitMap[updatedCommit.id] = {
          ...updatedCommit
        };
        this.buildCommitUpdate.next(!this.buildCommitUpdate.getValue());
      }
    });
  }

  public stopBuildCommitSubscription(): void {
    if (this.buildCommitSubscription) this.buildCommitSubscription.unsubscribe();
  }

  public startAgentTaskSubscription() {
    this.agentTaskSubscription = this.subscribeUpdatedAgentTask.subscribe().subscribe(({ data: { updatedAgentTask }, errors }) => {
      if (errors) {
        console.error(errors);
      } else if (updatedAgentTask) {
        this.agentTaskMap[updatedAgentTask.id] = {
          ...updatedAgentTask
        };
        this.agentTaskUpdate.next(!this.agentTaskUpdate.getValue());
      }
    });
  }

  public stopAgentTaskSubscription(): void {
    if (this.agentTaskSubscription) this.agentTaskSubscription.unsubscribe();
  }
}
