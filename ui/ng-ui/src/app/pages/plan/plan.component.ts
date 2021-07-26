import { ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { LaForgeGetBuildTreeQuery, LaForgeGetEnvironmentInfoQuery } from '@graphql';
import { Observable, Subscription } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

// import { Step } from '../.../../../models/plan.model';

@Component({
  selector: 'app-plan',
  templateUrl: './plan.component.html',
  styleUrls: ['./plan.component.scss']
})
export class PlanComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
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
    this.subheader.setShowEnvDropdown(true);

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
  }

  ngOnInit(): void {
    const sub1 = this.envService.getBuildTree().subscribe(() => {
      this.envService.initPlanStatuses();
      this.envService.initAgentStatuses();
      this.envService.initPlans();
    });
    this.unsubscribe.push(sub1);
  }

  ngOnDestroy(): void {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
  }

  envIsSelected(): boolean {
    return this.envService.getEnvironmentInfo().getValue() != null;
  }
}
