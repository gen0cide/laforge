import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { Routes, RouterModule } from '@angular/router';
import { LaforgePipesModule } from 'src/app/pipes/pipes.module';

import { ViewComponentsModule } from '../../components/view-components.module';

import { PlanComponent } from './plan.component';

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
    MatSelectModule,
    MatChipsModule,
    LaforgePipesModule
  ]
})
export class PlanModule {}
