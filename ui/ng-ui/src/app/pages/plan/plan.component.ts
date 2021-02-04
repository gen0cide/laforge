import { ChangeDetectorRef, Component, NgModule, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { EnvironmentInfo } from 'src/app/models/api.model';
import { Environment } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { EnvironmentService } from 'src/app/services/environment/environment.service';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';

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
  envs: EnvironmentInfo[];
  environment: Environment = null;
  apolloError: any = {};

  constructor(
    private api: ApiService,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    private envService: EnvironmentService
  ) {
    this.subheader.setTitle('Plan');
    this.subheader.setDescription('Plan an environment to build');
  }

  ngOnInit(): void {
    this.api.pullEnvironments().then((envs: EnvironmentInfo[]) => {
      this.envs = envs;
      this.cdRef.detectChanges();
    });
  }

  envIsSelected(): boolean {
    return this.envService.getCurrentEnv() != null && this.environment != null;
  }

  grabEnvironmentTree(changeEvent: MatSelectChange): void {
    this.envService.setCurrentEnv(this.envs.filter((e) => e.id === changeEvent.value)[0]);
    this.api.pullEnvTree(this.envService.getCurrentEnv().id).then(
      (env: Environment) => {
        this.environment = {
          ...env,
          build: {
            ...env.build,
            teams: [...env.build.teams]
              .sort((a, b) => a.teamNumber - b.teamNumber)
              .map((team) => ({
                ...team,
                provisionedNetworks: [...team.provisionedNetworks]
                  .sort((a, b) => {
                    if (a.name < b.name) return -1;
                    if (a.name > b.name) return 1;
                    return 0;
                  })
                  .map((network) => ({
                    ...network,
                    provisionedHosts: [...network.provisionedHosts].sort((a, b) => {
                      if (a.host.hostname < b.host.hostname) return -1;
                      if (a.host.hostname > b.host.hostname) return 1;
                      return 0;
                    })
                  }))
              }))
          }
        };
        this.cdRef.detectChanges();
      },
      (err) => {
        this.apolloError = err;
        this.cdRef.detectChanges();
      }
    );
  }
}
