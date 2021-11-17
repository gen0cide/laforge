import { Component, Input, OnInit, ChangeDetectorRef, OnDestroy } from '@angular/core';
import { LaForgeGetCurrentUserTasksQuery, LaForgeSubscribeUpdatedStatusSubscription } from '@graphql';
import { EnvironmentService } from '@services/environment/environment.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-task',
  templateUrl: './task.component.html',
  styleUrls: ['./task.component.scss']
})
export class TaskComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[];
  @Input() task: LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks'][0];
  taskStatus: LaForgeSubscribeUpdatedStatusSubscription['updatedStatus'];

  constructor(private envService: EnvironmentService, private cdRef: ChangeDetectorRef) {
    this.unsubscribe = [];
  }

  ngOnInit(): void {
    this.taskStatus = { ...this.task.ServerTaskToStatus };

    const sub = this.envService.statusUpdate.asObservable().subscribe(() => {
      this.checkTaskStatus();
      this.cdRef.detectChanges();
    });
    this.unsubscribe.push(sub);
  }

  ngOnDestroy(): void {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
  }

  checkTaskStatus(): void {
    if (!this.task.ServerTaskToStatus) return;
    const updatedStatus = this.envService.getStatus(this.task.ServerTaskToStatus.id);
    if (updatedStatus) {
      this.taskStatus = updatedStatus;
      this.cdRef.markForCheck();
    }
  }
}
