import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { BuildComponent } from './build.component';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';

import { ViewComponentsModule } from '../../components/view-components.module';
import { MatSelectModule } from '@angular/material/select';

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
    MatSelectModule
  ]
})
export class BuildModule {}
