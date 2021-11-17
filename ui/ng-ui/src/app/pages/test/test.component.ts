import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ProvisionStatus, Status } from 'src/app/models/common.model';
import { Environment, resolveEnvEnums } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';

@Component({
  selector: 'app-test',
  templateUrl: './test.component.html',
  styleUrls: ['./test.component.scss']
})
export class TestComponent implements OnInit {
  environment: Environment;
  failStatus: Status;
  completeStatus: Status;
  awaitingStatus: Status;
  inProgressStatus: Status;
  taintedStatus: Status;
  undefinedStatus: Status;

  constructor(private api: ApiService, private cdRef: ChangeDetectorRef) {
    this.completeStatus = {
      state: ProvisionStatus.COMPLETE,
      startedAt: Date.now().toString(),
      endedAt: Date.now().toString(),
      completed: true,
      failed: false,
      error: ''
    };
    this.failStatus = {
      state: ProvisionStatus.FAILED,
      startedAt: Date.now().toString(),
      endedAt: Date.now().toString(),
      completed: false,
      failed: true,
      error: ''
    };
    this.awaitingStatus = {
      state: ProvisionStatus.AWAITING,
      startedAt: Date.now().toString(),
      endedAt: Date.now().toString(),
      completed: false,
      failed: false,
      error: ''
    };
    this.inProgressStatus = {
      state: ProvisionStatus.INPROGRESS,
      startedAt: Date.now().toString(),
      endedAt: Date.now().toString(),
      completed: false,
      failed: false,
      error: ''
    };
    this.taintedStatus = {
      state: ProvisionStatus.TAINTED,
      startedAt: Date.now().toString(),
      endedAt: Date.now().toString(),
      completed: false,
      failed: true,
      error: ''
    };
    this.undefinedStatus = {
      state: ProvisionStatus.UNDEFINED,
      startedAt: Date.now().toString(),
      endedAt: Date.now().toString(),
      completed: false,
      failed: false,
      error: ''
    };
  }

  ngOnInit(): void {
    this.api.pullEnvironments().then((envs) => {
      this.api.getEnvironment(envs[0].id).subscribe((result) => {
        /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
        this.environment = resolveEnvEnums((result.data as any).environment) as Environment;
        this.cdRef.detectChanges();
      }, console.error);
    });
  }
}
