import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { InlineSVGModule } from 'ng-inline-svg';
import { SubheaderComponent } from './subheader/subheader.component';
import { NgbDropdownModule } from '@ng-bootstrap/ng-bootstrap';
import { DropdownMenusModule } from '../../content/dropdown-menus/dropdown-menus.module';

@NgModule({
  declarations: [
    SubheaderComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
    InlineSVGModule,
    NgbDropdownModule,
    DropdownMenusModule,
  ],
  exports: [SubheaderComponent],
})
export class SubheaderModule { }
