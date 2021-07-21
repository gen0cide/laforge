import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatTableModule } from '@angular/material/table';

import { MomentModule } from 'ngx-moment';

import { LaforgePipesModule } from '../pipes/pipes.module';

import { HostModalComponent } from './host-modal/host-modal.component';
import { HostComponent } from './host/host.component';
import { NetworkModalComponent } from './network-modal/network-modal.component';
import { NetworkComponent } from './network/network.component';
import { StepComponent } from './step/step.component';
import { TaskListComponent } from './task-list/task-list.component';
import { TaskComponent } from './task/task.component';
import { TeamComponent } from './team/team.component';

@NgModule({
  declarations: [
    NetworkComponent,
    HostComponent,
    TeamComponent,
    HostModalComponent,
    NetworkModalComponent,
    StepComponent,
    TaskListComponent,
    TaskComponent
  ],
  imports: [CommonModule, MatDialogModule, MatTableModule, MatButtonModule, MatCheckboxModule, MomentModule, LaforgePipesModule],
  exports: [NetworkComponent, HostComponent, TeamComponent, TaskListComponent, TaskComponent]
})
export class ViewComponentsModule {}
