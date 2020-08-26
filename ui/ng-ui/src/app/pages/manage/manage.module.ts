import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { ManageComponent } from './manage.component';
import { MatSliderModule } from '@angular/material/slider';
import { MatTableModule } from '@angular/material/table';

const routes: Routes = [
  {
    path: '',
    component: ManageComponent
  },
];

@NgModule({
  declarations: [ManageComponent],
  imports: [
    MatSliderModule,
    MatTableModule,
    CommonModule,
    RouterModule.forChild(routes)
  ],
})
export class ManageModule { }
