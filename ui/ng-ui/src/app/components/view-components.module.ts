import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatTableModule } from '@angular/material/table';

import { MomentModule } from 'ngx-moment';

import { LaforgePipesModule } from '../pipes/pipes.module';

import { DeleteBuildModalComponent } from './delete-build-modal/delete-build-modal.component';
import { HostModalComponent } from './host-modal/host-modal.component';
import { HostComponent } from './host/host.component';
import { NetworkModalComponent } from './network-modal/network-modal.component';
import { NetworkComponent } from './network/network.component';
import { StepModalComponent } from './step-modal/step-modal.component';
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
    DeleteBuildModalComponent,
    StepComponent,
    TaskListComponent,
    TaskComponent,
    StepModalComponent
  ],
  imports: [
    CommonModule,
    MatDialogModule,
    MatTableModule,
    MatButtonModule,
    MatCheckboxModule,
    MomentModule,
    LaforgePipesModule,
    MatInputModule,
    MatFormFieldModule
  ],
  exports: [NetworkComponent, HostComponent, TeamComponent, TaskListComponent, TaskComponent]
})
export class ViewComponentsModule {}
