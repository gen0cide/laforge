import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { ViewComponentsModule } from '@components/view-components.module';
import { NgbDropdownModule, NgbProgressbarModule } from '@ng-bootstrap/ng-bootstrap';
import { InlineSVGModule } from 'ng-inline-svg';

import { CoreModule } from '../_metronic/core';
import { ExtrasModule } from '../_metronic/partials/layout/extras/extras.module';
import { SubheaderModule } from '../_metronic/partials/layout/subheader/subheader.module';
import { TranslationModule } from '../modules/i18n/translation.module';

import { AsideDynamicComponent } from './_layout/components/aside-dynamic/aside-dynamic.component';
import { AsideComponent } from './_layout/components/aside/aside.component';
import { FooterComponent } from './_layout/components/footer/footer.component';
import { HeaderMobileComponent } from './_layout/components/header-mobile/header-mobile.component';
import { HeaderMenuComponent } from './_layout/components/header/header-menu/header-menu.component';
import { HeaderComponent } from './_layout/components/header/header.component';
import { LanguageSelectorComponent } from './_layout/components/topbar/language-selector/language-selector.component';
import { TopbarComponent } from './_layout/components/topbar/topbar.component';
import { ScriptsInitComponent } from './_layout/init/scipts-init/scripts-init.component';
import { LayoutComponent } from './_layout/layout.component';
import { PagesRoutingModule } from './pages-routing.module';

@NgModule({
  declarations: [
    LayoutComponent,
    ScriptsInitComponent,
    HeaderMobileComponent,
    AsideComponent,
    FooterComponent,
    HeaderComponent,
    HeaderMenuComponent,
    TopbarComponent,
    LanguageSelectorComponent,
    AsideDynamicComponent
  ],
  imports: [
    CommonModule,
    PagesRoutingModule,
    TranslationModule,
    InlineSVGModule,
    ExtrasModule,
    NgbDropdownModule,
    NgbProgressbarModule,
    CoreModule,
    SubheaderModule,
    ViewComponentsModule,
    MatSlideToggleModule
  ]
})
export class LayoutModule {}
