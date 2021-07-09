import { ChangeDetectorRef, Component, NgModule, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { LaForgeGetBuildTreeQuery, LaForgeGetEnvironmentInfoQuery } from '@graphql';
import { Observable } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { EnvironmentInfo } from 'src/app/models/api.model';
import { ID } from 'src/app/models/common.model';
import { Build, Environment } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

// import { Step } from '../.../../../models/plan.model';
import { PlanService } from '../../plan.service';
// import { chike } from '../../../data/sample-config';

// interface EnvConfig {
//   text : string
// }

@Component({
  selector: 'app-plan',
  templateUrl: './plan.component.html',
  styleUrls: ['./plan.component.scss']
})
export class PlanComponent implements OnInit {
  environment: Observable<LaForgeGetEnvironmentInfoQuery['environment']>;
  build: Observable<LaForgeGetBuildTreeQuery['build']>;
  apolloError: any = {};

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    public envService: EnvironmentService
  ) {
    this.subheader.setTitle('Plan');
    this.subheader.setDescription('Plan an environment to build');

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
  }

  ngOnInit(): void {}

  envIsSelected(): boolean {
    return this.envService.getEnvironmentInfo().getValue() != null;
  }
}
