import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';

import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { Routes, RouterModule } from '@angular/router';

import { ViewComponentsModule } from '../../components/view-components.module';

import { AccountComponent } from './account.component';

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
    MatButtonModule,
    ReactiveFormsModule,
    FormsModule
  ]
})
export class AccountModule {}
