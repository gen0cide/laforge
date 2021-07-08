import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { interval, Observable } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { Build, Environment } from 'src/app/models/environment.model';
import { EnvironmentService } from 'src/app/services/environment/environment.service';

@Component({
  selector: 'app-build',
  templateUrl: './build.component.html',
  styleUrls: ['./build.component.scss']
})
export class BuildComponent implements OnInit {
  environment: Observable<Environment>;
  build: Observable<Build>;

  constructor(private subheader: SubheaderService, public envService: EnvironmentService, private cdRef: ChangeDetectorRef) {
    this.subheader.setTitle('Build');
    this.subheader.setDescription('Build a planned environment');

    this.environment = this.envService.getCurrentEnv().asObservable();
    this.build = this.envService.getCurrentBuild().asObservable();
  }

  ngOnInit(): void {
    interval(10000).subscribe(() => {
      this.envService.updatePlanStatuses();
      this.cdRef.detectChanges();
    });
  }

  envIsSelected(): boolean {
    return this.envService.getCurrentEnv().getValue() != null;
  }
}
