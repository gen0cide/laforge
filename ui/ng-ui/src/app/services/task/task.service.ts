import { Injectable } from '@angular/core';
import { LaForgeGetCurrentUserTasksGQL, LaForgeGetCurrentUserTasksQuery, LaForgeSubscribeUpdatedServerTaskGQL } from '@graphql';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TaskService {
  private tasks: BehaviorSubject<LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks']>;
  private taskNotification: BehaviorSubject<boolean>;

  constructor(
    private getCurrentUserTasks: LaForgeGetCurrentUserTasksGQL,
    private subscribeUpdatedServerTask: LaForgeSubscribeUpdatedServerTaskGQL
  ) {
    this.tasks = new BehaviorSubject([]);
    this.taskNotification = new BehaviorSubject(false);

    this.initTasks();
    this.startTaskSubscription();
  }

  public getTasks(): BehaviorSubject<LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks']> {
    return this.tasks;
  }

  public getTaskNotification(): BehaviorSubject<boolean> {
    return this.taskNotification;
  }

  public notifyUser(): void {
    this.taskNotification.next(true);
  }

  public clearUserNotification(): void {
    this.taskNotification.next(false);
  }

  private initTasks(): void {
    this.getCurrentUserTasks
      .fetch()
      .toPromise()
      .then(({ data }) => {
        if (data) {
          this.tasks.next(
            [...data.getCurrentUserTasks].sort((a, b) => new Date(b.start_time).valueOf() - new Date(a.start_time).valueOf())
          );
        }
      });
  }

  private startTaskSubscription(): void {
    this.subscribeUpdatedServerTask.subscribe().subscribe(({ data: { updatedServerTask }, errors }) => {
      if (errors) {
        console.error(errors);
      } else if (updatedServerTask) {
        const currentTasks = this.tasks.getValue();
        const taskIndex = currentTasks.map((t) => t.id).indexOf(updatedServerTask.id);
        if (taskIndex >= 0) {
          currentTasks[taskIndex] = {
            ...updatedServerTask
          };
        } else {
          currentTasks.push({ ...updatedServerTask });
        }
        this.tasks.next(currentTasks.sort((a, b) => new Date(b.start_time).valueOf() - new Date(a.start_time).valueOf()));
        this.notifyUser();
      }
    });
  }
}
