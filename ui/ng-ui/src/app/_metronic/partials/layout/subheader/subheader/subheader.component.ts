import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { BreadcrumbItemModel } from '../_models/breadcrumb-item.model';
import { LayoutService } from '../../../../core';
import { SubheaderService } from '../_services/subheader.service';

import { PlanService } from '../../../../../plan.service';
import { MatSelectChange } from '@angular/material/select';
import { EnvironmentService } from 'src/app/services/environment/environment.service';
import { Environment } from 'src/app/models/environment.model';

interface Branch {
  name: string;
  hash: string;
}

@Component({
  selector: 'app-subheader',
  templateUrl: './subheader.component.html',
  styleUrls: ['./subheader.component.scss']
})
export class SubheaderComponent implements OnInit {
  subheaderCSSClasses = '';
  subheaderContainerCSSClasses = '';
  subheaderMobileToggle = false;
  subheaderDisplayDesc = false;
  subheaderDisplayDaterangepicker = false;
  title$: Observable<string>;
  breadcrumbs$: Observable<BreadcrumbItemModel[]>;
  description$: Observable<string>;
  branches: Branch[] = [
    { name: 'Bradley', hash: '98y3if' },
    { name: 'Lucas', hash: '32a7fh' }
  ];
  environment: Environment;
  envIsLoading: Observable<boolean>;

  constructor(
    private layout: LayoutService,
    private subheader: SubheaderService,
    private planService: PlanService,
    public envService: EnvironmentService,
    private cdRef: ChangeDetectorRef
  ) {
    this.title$ = this.subheader.titleSubject.asObservable();
    this.breadcrumbs$ = this.subheader.breadCrumbsSubject.asObservable();
    this.description$ = this.subheader.descriptionSubject.asObservable();

    this.envService.getEnvironments().subscribe(() => {
      this.cdRef.markForCheck();
    });

    this.envService
      .getCurrentEnv()
      .asObservable()
      .subscribe((env) => (this.environment = env));
    this.envIsLoading = this.envService.envIsLoading.asObservable();
  }

  ngOnInit() {
    this.subheaderCSSClasses = this.layout.getStringCSSClasses('subheader');
    this.subheaderContainerCSSClasses = this.layout.getStringCSSClasses('subheader_container');
    this.subheaderMobileToggle = this.layout.getProp('subheader.mobileToggle');
    this.subheaderDisplayDesc = this.layout.getProp('subheader.displayDesc');
    this.subheaderDisplayDaterangepicker = this.layout.getProp('subheader.displayDaterangepicker');
  }

  onBranchSelect(changeEvent: MatSelectChange) {
    // console.log(changeEvent);
    this.planService.getPlan(changeEvent.value);
  }

  selectEnvironment(event: MatSelectChange) {
    this.envService.setCurrentEnv(event.value);
    this.cdRef.detectChanges();
  }
}