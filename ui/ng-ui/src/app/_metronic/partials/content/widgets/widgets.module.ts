import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgApexchartsModule } from 'ng-apexcharts';
import { InlineSVGModule } from 'ng-inline-svg';
import { MixedWidget1Component } from './mixed/mixed-widget1/mixed-widget1.component';
import { ListsWidget9Component } from './lists/lists-widget9/lists-widget9.component';
import { StatsWidget11Component } from './stats/stats-widget11/stats-widget11.component';
import { StatsWidget12Component } from './stats/stats-widget12/stats-widget12.component';
import { ListsWidget1Component } from './lists/lists-widget1/lists-widget1.component';
import { ListsWidget3Component } from './lists/lists-widget3/lists-widget3.component';
import { ListsWidget4Component } from './lists/lists-widget4/lists-widget4.component';
import { ListsWidget8Component } from './lists/lists-widget8/lists-widget8.component';
import { MixedWidget14Component } from './mixed/mixed-widget14/mixed-widget14.component';
import { AdvanceTablesWidget4Component } from './advance-tables/advance-tables-widget4/advance-tables-widget4.component';
import { AdvanceTablesWidget2Component } from './advance-tables/advance-tables-widget2/advance-tables-widget2.component';
import { DropdownMenusModule } from '../dropdown-menus/dropdown-menus.module';
import { NgbDropdownModule } from '@ng-bootstrap/ng-bootstrap';
import { Widget3DropdownComponent } from './lists/lists-widget3/widget3-dropdown/widget3-dropdown.component';
import { Widget4DropdownComponent } from './lists/lists-widget4/widget4-dropdown/widget4-dropdown.component';

@NgModule({
  declarations: [
    MixedWidget1Component,
    ListsWidget9Component,
    StatsWidget11Component,
    StatsWidget12Component,
    ListsWidget1Component,
    ListsWidget3Component,
    ListsWidget4Component,
    ListsWidget8Component,
    MixedWidget14Component,
    AdvanceTablesWidget4Component,
    AdvanceTablesWidget2Component,
    Widget3DropdownComponent,
    Widget4DropdownComponent,
  ],
  imports: [
    CommonModule,
    DropdownMenusModule,
    InlineSVGModule,
    NgApexchartsModule,
    NgbDropdownModule,
  ],
  exports: [
    MixedWidget1Component,
    ListsWidget9Component,
    StatsWidget11Component,
    StatsWidget12Component,
    ListsWidget1Component,
    ListsWidget3Component,
    ListsWidget4Component,
    ListsWidget8Component,
    MixedWidget14Component,
    AdvanceTablesWidget4Component,
    AdvanceTablesWidget2Component,
  ],
})
export class WidgetsModule {}
