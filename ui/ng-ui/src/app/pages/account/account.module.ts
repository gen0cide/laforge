import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { AccountComponent } from './account.component';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';

import { ViewComponentsModule } from '../../components/view-components.module';
import { MatSelectModule } from '@angular/material/select';
import { MatChipsModule } from '@angular/material/chips';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatDividerModule } from '@angular/material/divider';

const routes: Routes = [
  {
    path: '',
    component: AccountComponent
  }
];

@NgModule({
  declarations: [AccountComponent],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    ViewComponentsModule,
    MatCardModule,
    MatTableModule,
    MatSelectModule,
    MatChipsModule,
    MatInputModule,
    MatFormFieldModule,
    MatDividerModule,
    MatButtonModule
  ]
})
export class AccountModule {}
