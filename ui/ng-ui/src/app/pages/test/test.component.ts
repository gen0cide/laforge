import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { ProvisionStatus } from 'src/app/models/common.model';

@Component({
  selector: 'app-test',
  templateUrl: './test.component.html',
  styleUrls: ['./test.component.scss']
})
export class TestComponent implements OnInit {
  environment: Environment;
  failStatus: ProvisionStatus = ProvisionStatus.ProvStatusFailed;
  completeStatus: ProvisionStatus = ProvisionStatus.ProvStatusComplete;
  awaitingStatus: ProvisionStatus = ProvisionStatus.ProvStatusAwaiting;
  inProgressStatus: ProvisionStatus = ProvisionStatus.ProvStatusInProgress;
  taintedStatus: ProvisionStatus = ProvisionStatus.ProvStatusTainted;
  undefinedStatus: ProvisionStatus = ProvisionStatus.ProvStatusUndefined;

  constructor(private api: ApiService, private cdRef: ChangeDetectorRef) {}

  ngOnInit(): void {
    this.api.getEnvironment('a3f73ee0-da71-4aa6-9280-18ad1a1a8d16').subscribe((result) => {
      /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
      this.environment = resolveStatuses((result.data as any).environment) as Environment;
      this.cdRef.detectChanges();
    });
  }
}
