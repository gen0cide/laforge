import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { MatSliderModule } from '@angular/material/slider';

const routes: Routes = [
  {
    path: '',
    component: DashboardComponent
  },
];

@NgModule({
  declarations: [DashboardComponent],
  imports: [
    MatSliderModule,
    CommonModule,
    RouterModule.forChild(routes),
  ],
})
export class DashboardModule { }
