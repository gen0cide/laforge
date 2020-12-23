import { Component, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { Location } from '@angular/common';
import { LayoutService, DynamicAsideMenuService } from '../../../../_metronic/core';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-aside-dynamic',
  templateUrl: './aside-dynamic.component.html',
  styleUrls: ['./aside-dynamic.component.scss']
})
export class AsideDynamicComponent implements OnInit, OnDestroy {
  menuConfig: any;
  subscriptions: Subscription[] = [];

  disableAsideSelfDisplay: boolean;
  headerLogo: string;
  brandSkin: string;
  ulCSSClasses: string;
  location: Location;
  asideMenuHTMLAttributes: any = {};
  asideMenuCSSClasses: string;
  brandClasses: string;
  asideMenuScroll = 1;
  asideSelfMinimizeToggle = false;

  constructor(
    private layout: LayoutService,
    private loc: Location,
    private menu: DynamicAsideMenuService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    // load view settings
    this.disableAsideSelfDisplay = this.layout.getProp('aside.self.display') === false;
    this.brandSkin = this.layout.getProp('brand.self.theme');
    this.headerLogo = this.getLogo();
    this.ulCSSClasses = this.layout.getProp('aside_menu_nav');
    this.asideMenuCSSClasses = this.layout.getStringCSSClasses('aside_menu');
    this.asideMenuHTMLAttributes = this.layout.getHTMLAttributes('aside_menu');
    this.brandClasses = this.layout.getProp('brand');
    this.asideSelfMinimizeToggle = this.layout.getProp('aside.self.minimize.toggle');
    this.asideMenuScroll = this.layout.getProp('aside.menu.scroll') ? 1 : 0;
    this.location = this.loc;

    // load menu
    const subscr = this.menu.menuConfig$.subscribe((res) => {
      this.menuConfig = res;
      this.cdr.detectChanges();
    });
    this.subscriptions.push(subscr);
  }

  private getLogo() {
    if (this.brandSkin === 'light') {
      return '/assets/media/logos/logo-dark.png';
    } else {
      return '/assets/media/logos/logo-light.png';
    }
  }

  ngOnDestroy() {
    this.subscriptions.forEach((sb) => sb.unsubscribe());
  }
}
