import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { ManageComponent } from './manage.component';
import { MatCardModule } from '@angular/material/card';

import { ViewComponentsModule } from '../../components/view-components.module';

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
    ViewComponentsModule
  ]
})
export class ManageModule {}
