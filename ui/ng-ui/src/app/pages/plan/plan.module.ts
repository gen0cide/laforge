import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { PlanComponent } from './plan.component';
import { MatCardModule } from '@angular/material/card';

import { ViewComponentsModule } from '../../components/view-components.module';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';

const routes: Routes = [
  {
    path: '',
    component: PlanComponent
  }
];

@NgModule({
  declarations: [PlanComponent],
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
export class PlanModule {}
