import { Injectable } from '@angular/core';
import {
  LaForgeGetBuildCommitsGQL,
  LaForgeGetBuildCommitsQuery,
  LaForgeGetBuildPlansGQL,
  LaForgeGetBuildPlansQuery,
  LaForgeGetBuildTreeGQL,
  LaForgeGetBuildTreeQuery,
  LaForgeGetEnvironmentGQL,
  LaForgeGetEnvironmentInfoQuery,
  LaForgeGetEnvironmentsGQL,
  LaForgeGetEnvironmentsQuery,
  LaForgePullAgentStatusesGQL,
  LaForgePullAgentStatusesQuery,
  LaForgePullPlanStatusesGQL,
  LaForgePullPlanStatusesQuery
} from '@graphql';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(
    private getEnvironments: LaForgeGetEnvironmentsGQL,
    private pullPlanStatuses: LaForgePullPlanStatusesGQL,
    private pullAgentStatuses: LaForgePullAgentStatusesGQL,
    private getEnvironmentInfoGQL: LaForgeGetEnvironmentGQL,
    private getBuildTreeGQL: LaForgeGetBuildTreeGQL,
    private getBuildPlansGQL: LaForgeGetBuildPlansGQL,
    private getBuildCommitsGQL: LaForgeGetBuildCommitsGQL
  ) {}

  /**
   * Pulls status objects for all plans on a build
   * @param buildId The build ID that contains plans
   * @returns All plan objects relating to a build
   */
  public pullAllPlanStatuses(buildId: string): Promise<LaForgePullPlanStatusesQuery['build']['buildToPlan']> {
    return new Promise((resolve, reject) => {
      this.pullPlanStatuses
        .fetch({
          buildId
        })
        .toPromise()
        .then(({ data, error, errors }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.build.buildToPlan);
        });
    });
  }

  /**
   * Pulls the build tree (with all branches only having ID's) and its contained agent statuses
   * @param buildId The build ID that agents relate to
   * @returns The build tree with only agents as full objects
   */
  public pullAllAgentStatuses(buildId: string): Promise<LaForgePullAgentStatusesQuery['build']> {
    return new Promise((resolve, reject) => {
      this.pullAgentStatuses
        .fetch({
          buildId
        })
        .toPromise()
        .then(({ data, error, errors }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.build);
        });
    });
  }

  /**
   * Pulls an environment from the API once, without exposing a subscription or observable
   * @param id The Environment ID of the environment
   */
  public async pullEnvironments(): Promise<LaForgeGetEnvironmentsQuery['environments']> {
    return new Promise((resolve, reject) => {
      this.getEnvironments
        .fetch()
        .toPromise()
        .then(({ data, error, errors }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.environments);
        });
    });
  }

  public async pullEnvironmentInfo(envId: string): Promise<LaForgeGetEnvironmentInfoQuery['environment']> {
    return new Promise((resolve, reject) => {
      this.getEnvironmentInfoGQL
        .fetch({
          envId: envId
        })
        .toPromise()
        .then(({ data, error, errors }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.environment);
        }, reject);
    });
  }

  public async pullBuildTree(buildId: string): Promise<LaForgeGetBuildTreeQuery['build']> {
    return new Promise((resolve, reject) => {
      this.getBuildTreeGQL
        .fetch({
          buildId: buildId as string
        })
        .toPromise()
        .then(({ data, error, errors }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.build);
        }, reject);
    });
  }

  public async pullBuildPlans(buildId: string): Promise<LaForgeGetBuildPlansQuery['build']> {
    return new Promise((resolve, reject) => {
      this.getBuildPlansGQL
        .fetch({
          buildId
        })
        .toPromise()
        .then(({ data, errors, error }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.build);
        }, reject);
    });
  }

  public async pullBuildCommits(buildId: string): Promise<LaForgeGetBuildCommitsQuery['build']['BuildToBuildCommits']> {
    return new Promise((resolve, reject) => {
      this.getBuildCommitsGQL
        .fetch({
          buildId
        })
        .toPromise()
        .then(({ data, errors, error }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.build.BuildToBuildCommits);
        }, reject);
    });
  }
}
