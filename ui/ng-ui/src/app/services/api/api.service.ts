import { Injectable } from '@angular/core';
import {
  LaForgeCreateBuildGQL,
  LaForgeCreateBuildMutation,
  LaForgeCreateEnvironmentFromGitGQL,
  LaForgeCreateEnvironmentFromGitMutation,
  LaForgeCreateUserGQL,
  LaForgeCreateUserMutation,
  LaForgeGetAllAgentStatusesGQL,
  LaForgeGetAllAgentStatusesQuery,
  LaForgeGetAllPlanStatusesGQL,
  LaForgeGetAllPlanStatusesQuery,
  LaForgeGetBuildCommitsGQL,
  LaForgeGetBuildCommitsQuery,
  LaForgeGetBuildPlansGQL,
  LaForgeGetBuildPlansQuery,
  LaForgeGetBuildTreeGQL,
  LaForgeGetBuildTreeQuery,
  LaForgeGetEnvironmentInfoGQL,
  LaForgeGetEnvironmentInfoQuery,
  LaForgeGetEnvironmentsGQL,
  LaForgeGetEnvironmentsQuery,
  LaForgeGetUserListGQL,
  LaForgeGetUserListQuery,
  LaForgeModifyCurrentUserGQL,
  LaForgeModifyCurrentUserMutation,
  LaForgeProviderType,
  LaForgeRoleLevel,
  LaForgeUpdateUserGQL,
  LaForgeUpdateUserMutation,
  LaForgeUpdateEnviromentViaPullGQL,
  LaForgeUpdateEnviromentViaPullMutation
} from '@graphql';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(
    private getEnvironments: LaForgeGetEnvironmentsGQL,
    // private pullPlanStatuses: LaForgePullPlanStatusesGQL,
    // private pullAgentStatuses: LaForgePullAgentStatusesGQL,
    private getAllAgentStatuses: LaForgeGetAllAgentStatusesGQL,
    private getAllPlanStatuses: LaForgeGetAllPlanStatusesGQL,
    private getEnvironmentInfoGQL: LaForgeGetEnvironmentInfoGQL,
    private getBuildTreeGQL: LaForgeGetBuildTreeGQL,
    private getBuildPlansGQL: LaForgeGetBuildPlansGQL,
    private getBuildCommitsGQL: LaForgeGetBuildCommitsGQL,
    private createBuildGQL: LaForgeCreateBuildGQL,
    private modifyCurrentUserGQL: LaForgeModifyCurrentUserGQL,
    private createEnvironmentFromGitGQL: LaForgeCreateEnvironmentFromGitGQL,
    private getUserListGQL: LaForgeGetUserListGQL,
    private updateUserGQL: LaForgeUpdateUserGQL,
    private createUserGQL: LaForgeCreateUserGQL,
    private updateEnviromentViaPullGQL: LaForgeUpdateEnviromentViaPullGQL
  ) {}

  /**
   * Pulls status objects for all plans on a build
   * @param buildId The build ID that contains plans
   * @returns All plan objects relating to a build
   */
  public pullAllPlanStatuses(buildId: string, count: number, offset: number): Promise<LaForgeGetAllPlanStatusesQuery['getAllPlanStatus']> {
    return new Promise((resolve, reject) => {
      this.getAllPlanStatuses
        .fetch({
          buildId,
          count,
          offset
        })
        .toPromise()
        .then(({ data, error, errors }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.getAllPlanStatus);
        });
    });
  }

  /**
   * Pulls the build tree (with all branches only having ID's) and its contained agent statuses
   * @param buildId The build ID that agents relate to
   * @returns The build tree with only agents as full objects
   */
  public pullAllAgentStatuses(
    buildId: string,
    count: number,
    offset: number
  ): Promise<LaForgeGetAllAgentStatusesQuery['getAllAgentStatus']> {
    return new Promise((resolve, reject) => {
      this.getAllAgentStatuses
        .fetch({
          buildId,
          count,
          offset
        })
        .toPromise()
        .then(({ data, error, errors }) => {
          if (error) {
            return reject(error);
          } else if (errors) {
            return reject(errors);
          }
          resolve(data.getAllAgentStatus);
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

  public async createBuild(envId: string): Promise<LaForgeCreateBuildMutation['createBuild']> {
    return new Promise((resolve, reject) => {
      this.createBuildGQL
        .mutate({
          envId
        })
        .toPromise()
        .then(({ data, errors }) => {
          if (errors) {
            return reject(errors);
          } else if (data.createBuild) {
            return resolve(data.createBuild);
          }
          reject(new Error('unknown error occurred while creating build'));
        }, reject);
    });
  }

  public async updateEnvFromGit(envId: string): Promise<LaForgeUpdateEnviromentViaPullMutation['updateEnviromentViaPull']> {
    return new Promise((resolve, reject) => {
      this.updateEnviromentViaPullGQL
        .mutate({
          envId
        })
        .toPromise()
        .then(({ data, errors }) => {
          if (errors) {
            return reject(errors);
          } else if (data) {
            return resolve(data.updateEnviromentViaPull);
          }
          reject(new Error('unknown error occurred while updating enviroment'));
        }, reject);
    });
  }

  public async updateAuthUser(updateAuthUserInput: {
    firstName?: string;
    lastName?: string;
    email?: string;
    phone?: string;
    company?: string;
    occupation?: string;
  }): Promise<LaForgeModifyCurrentUserMutation['modifySelfUserInfo']> {
    return new Promise((resolve, reject) => {
      this.modifyCurrentUserGQL
        .mutate({
          ...updateAuthUserInput
        })
        .toPromise()
        .then(({ data, errors }) => {
          if (errors) {
            return reject(errors);
          } else if (data.modifySelfUserInfo) {
            return resolve(data.modifySelfUserInfo);
          }
          reject(new Error('unknown error occurred while updating current user'));
        }, reject);
    });
  }

  public async createEnvFromGit(createEnvFromGitInput: {
    repoURL: string;
    repoName: string;
    branchName: string;
    envFilePath: string;
  }): Promise<LaForgeCreateEnvironmentFromGitMutation['createEnviromentFromRepo']> {
    return new Promise((resolve, reject) => {
      this.createEnvironmentFromGitGQL
        .mutate({
          ...createEnvFromGitInput
        })
        .toPromise()
        .then(({ data, errors }) => {
          if (errors) {
            return reject(errors);
          } else if (data.createEnviromentFromRepo) {
            return resolve(data.createEnviromentFromRepo);
          }
          reject(new Error('unknown error occurred while cloning env from git'));
        }, reject);
    });
  }

  public async getAllUsers(): Promise<LaForgeGetUserListQuery['getUserList']> {
    return new Promise((resolve, reject) => {
      this.getUserListGQL
        .fetch()
        .toPromise()
        .then(({ data, errors }) => {
          if (errors) {
            return reject(errors);
          } else if (data.getUserList) {
            return resolve(data.getUserList);
          }
          reject(new Error('unknown error occurred while getting user list'));
        }, reject);
    });
  }

  public async modifyUser(
    userId: string,
    input: {
      firstName?: string;
      lastName?: string;
      email: string;
      phone?: string;
      company?: string;
      occupation?: string;
      role: LaForgeRoleLevel;
      provider: LaForgeProviderType;
    }
  ): Promise<LaForgeUpdateUserMutation['modifyAdminUserInfo']> {
    return new Promise((resolve, reject) => {
      this.updateUserGQL
        .mutate({
          userId,
          ...input
        })
        .toPromise()
        .then(({ data, errors }) => {
          if (errors) {
            return reject(errors);
          } else if (data.modifyAdminUserInfo) {
            return resolve(data.modifyAdminUserInfo);
          }
          reject(new Error('unknown error occurred while updating user'));
        }, reject);
    });
  }

  public async createUser(input: {
    username: string;
    password: string;
    role: LaForgeRoleLevel;
    provider: LaForgeProviderType;
  }): Promise<LaForgeCreateUserMutation['createUser']> {
    return new Promise((resolve, reject) => {
      this.createUserGQL
        .mutate({
          ...input
        })
        .toPromise()
        .then(({ data, errors }) => {
          if (errors) {
            return reject(errors);
          } else if (data.createUser) {
            return resolve(data.createUser);
          }
          reject(new Error('unknown error occurred while creating user'));
        }, reject);
    });
  }
}
