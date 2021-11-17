import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { Routes, RouterModule } from '@angular/router';
import { LaforgePipesModule } from 'src/app/pipes/pipes.module';

import { ViewComponentsModule } from '../../components/view-components.module';

import { BuildComponent } from './build.component';

const routes: Routes = [
  {
    path: '',
    component: BuildComponent
  }
];

@NgModule({
  declarations: [BuildComponent],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    ViewComponentsModule,
    MatCardModule,
    MatTableModule,
    MatButtonModule,
    MatSelectModule,
    LaforgePipesModule
  ]
})
export class BuildModule {}
