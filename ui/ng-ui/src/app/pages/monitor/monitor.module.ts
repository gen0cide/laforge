import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { Routes, RouterModule } from '@angular/router';

import { ViewComponentsModule } from '../../components/view-components.module';

import { MonitorComponent } from './monitor.component';
// import { TeamComponent } from 'src/app/components/view-components.module';

const routes: Routes = [
  {
    path: '',
    component: MonitorComponent
  }
];

@NgModule({
  declarations: [MonitorComponent],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    MatCardModule,
    ViewComponentsModule,
    MatTableModule,
    MatButtonModule,
    MatSelectModule
  ]
})
export class MonitorModule {}
