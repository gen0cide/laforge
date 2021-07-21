import { ChangeDetectorRef, Component, OnInit, OnDestroy } from '@angular/core';
import { LaForgeGetBuildTreeQuery, LaForgeGetEnvironmentInfoQuery } from '@graphql';
import { Observable, Subscription } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

@Component({
  selector: 'app-build',
  templateUrl: './build.component.html',
  styleUrls: ['./build.component.scss']
})
export class BuildComponent implements OnInit, OnDestroy {
  private unsubscribe: Subscription[] = [];
  environment: Observable<LaForgeGetEnvironmentInfoQuery['environment']>;
  build: Observable<LaForgeGetBuildTreeQuery['build']>;

  constructor(private subheader: SubheaderService, public envService: EnvironmentService, private cdRef: ChangeDetectorRef) {
    this.subheader.setTitle('Build');
    this.subheader.setDescription('Build a planned environment');

    this.environment = this.envService.getEnvironmentInfo().asObservable();
    this.build = this.envService.getBuildTree().asObservable();
  }

  ngOnInit(): void {
    // interval(10000).subscribe(() => {
    //   this.envService.updatePlanStatuses();
    //   this.cdRef.detectChanges();
    // });
    const subr = this.envService.getBuildTree().subscribe(() => {
      this.envService.initPlanStatuses();
      this.envService.initAgentStatuses();
    });
    this.unsubscribe.push(subr);
  }

  ngOnDestroy(): void {
    this.unsubscribe.forEach((sub) => sub.unsubscribe());
  }

  envIsSelected(): boolean {
    return this.envService.getEnvironmentInfo().getValue() != null;
  }
}
