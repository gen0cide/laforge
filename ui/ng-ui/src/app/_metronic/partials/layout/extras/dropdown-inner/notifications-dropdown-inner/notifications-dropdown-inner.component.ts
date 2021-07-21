import { Component, OnInit } from '@angular/core';
import { LaForgeGetCurrentUserTasksQuery, LaForgeProvisionStatus, LaForgeServerTaskType } from '@graphql';
import { TaskService } from '@services/task/task.service';
import { LayoutService } from '../../../../../core';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-notifications-dropdown-inner',
  templateUrl: './notifications-dropdown-inner.component.html',
  styleUrls: ['./notifications-dropdown-inner.component.scss']
})
export class NotificationsDropdownInnerComponent implements OnInit {
  extrasNotificationsDropdownStyle: 'light' | 'dark' = 'dark';
  activeTabId: 'topbar_notifications_notifications' | 'topbar_notifications_events' | 'topbar_notifications_logs' =
    'topbar_notifications_notifications';
  tasks: Observable<LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks']>;

  constructor(private layout: LayoutService, private taskService: TaskService) {}

  ngOnInit(): void {
    this.extrasNotificationsDropdownStyle = this.layout.getProp('extras.notifications.dropdown.style');
    this.tasks = this.taskService.getTasks().asObservable();
  }

  setActiveTabId(tabId) {
    this.activeTabId = tabId;
  }

  getActiveCSSClasses(tabId) {
    if (tabId !== this.activeTabId) {
      return '';
    }
    return 'active show';
  }

  getMessageSentiment(status: LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks'][0]['ServerTaskToStatus']): string {
    switch (status.state) {
      case LaForgeProvisionStatus.Complete:
        return 'has completed';
      case LaForgeProvisionStatus.Failed:
        return 'has failed';
      case LaForgeProvisionStatus.Inprogress:
        return 'has been queued';
      default:
        return '';
    }
  }

  getMessageSubject(task: LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks'][0]): string {
    switch (task.type) {
      case LaForgeServerTaskType.Createbuild:
        if (task.ServerTaskToBuild) return `Create build '${task.ServerTaskToEnvironment?.name} v${task.ServerTaskToBuild?.revision}'`;
        else return `Creating build of '${task.ServerTaskToEnvironment?.name}'`;
      case LaForgeServerTaskType.Deletebuild:
        return `Deleting build '${task.ServerTaskToEnvironment?.name} v${task.ServerTaskToBuild?.revision}'`;
      case LaForgeServerTaskType.Loadenv:
        if (task.ServerTaskToEnvironment) return `Load environment '${task.ServerTaskToEnvironment?.name}'`;
        else return `Loading environment`;
      case LaForgeServerTaskType.Rebuild:
        return `Rebuilding of '${task.ServerTaskToEnvironment?.name} v${task.ServerTaskToBuild?.revision}'`;
      case LaForgeServerTaskType.Renderfiles:
        if (task.ServerTaskToBuild && task.ServerTaskToEnvironment)
          return `Render files for '${task.ServerTaskToEnvironment?.name} v${task.ServerTaskToBuild?.revision}'`;
        else if (task.ServerTaskToEnvironment) return `Render files for '${task.ServerTaskToEnvironment?.name}'`;
        else return `Rendering files`;
      default:
        return 'Unknown task';
    }
  }

  getMessage(task: LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks'][0]): string {
    return `${this.getMessageSubject(task)} ${this.getMessageSentiment(task.ServerTaskToStatus)}.`;
  }

  getIconClass(task: LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks'][0]): string {
    // fa-hammer text-success
    return `fa-${this.getIcon(task.type)} text-${this.getColor(task.ServerTaskToStatus)}`;
  }

  getColor(status: LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks'][0]['ServerTaskToStatus']): string {
    // fa-hammer text-success
    switch (status.state) {
      case LaForgeProvisionStatus.Complete:
        return 'success';
      case LaForgeProvisionStatus.Inprogress:
        return 'info';
      case LaForgeProvisionStatus.Failed:
        return 'danger';
      default:
        return 'dark';
    }
  }

  getIcon(type: LaForgeGetCurrentUserTasksQuery['getCurrentUserTasks'][0]['type']): string {
    // fa-hammer text-success
    switch (type) {
      case LaForgeServerTaskType.Createbuild:
        return 'hammer';
      case LaForgeServerTaskType.Deletebuild:
        return 'trash-alt';
      case LaForgeServerTaskType.Loadenv:
        return 'folder-tree';
      case LaForgeServerTaskType.Rebuild:
        return 'redo-alt';
      case LaForgeServerTaskType.Renderfiles:
        return 'print';
      default:
        return 'question';
    }
  }
}
