import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { MonitorComponent } from './monitor.component';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';

import { ViewComponentsModule } from '../../components/view-components.module';

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
