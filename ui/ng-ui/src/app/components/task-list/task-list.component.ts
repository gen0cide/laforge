import { Component, OnInit } from '@angular/core';
import { LaForgeGetCurrentUserTasksQuery } from '@graphql';
import { TaskService } from '@services/task/task.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-task-list',
  templateUrl: './task-list.component.html',
  styleUrls: ['./task-list.component.scss']
})
export class TaskListComponent implements OnInit {
  tasks: Observable<LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks']>;
  constructor(private taskService: TaskService) {}

  ngOnInit(): void {
    this.tasks = this.taskService.getTasks().asObservable();
  }
}
