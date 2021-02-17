import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { ManageComponent } from './manage.component';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';

import { ViewComponentsModule } from '../../components/view-components.module';
import { MatSelectModule } from '@angular/material/select';

const routes: Routes = [
  {
    path: '',
    component: ManageComponent
  }
];

@NgModule({
  declarations: [ManageComponent],
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
export class ManageModule {}
