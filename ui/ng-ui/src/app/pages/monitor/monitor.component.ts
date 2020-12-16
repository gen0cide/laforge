import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { ApolloQueryResult } from '@apollo/client/core';
import { QueryRef } from 'apollo-angular';
import { EmptyObject } from 'apollo-angular/types';
import { Observable, Subscription } from 'rxjs';
import { updateAgentStatuses } from 'src/app/models/agent.model';
import { AgentStatusQueryResult } from 'src/app/models/common.model';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
// import { ProvisionedNetwork } from 'src/app/models/network.model';

// import { corp_network_provisioned } from '../../../data/corp';
// import { bradley } from 'src/data/sample-config';
@Component({
  selector: 'app-manage',
  templateUrl: './monitor.component.html',
  styleUrls: ['./monitor.component.scss']
})
export class MonitorComponent implements OnInit {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  environment: Environment = null;
  loaded = false;
  displayedColumns: string[] = ['TeamCount', 'AdminCIDRs', 'ExposedVDIPorts', 'maintainer'];
  // agentStatusQueryRef: QueryRef<any, EmptyObject>;
  environmentSubscription: Subscription;
  // agentStatusSubscription: Subscription;
  agentPollingInterval: NodeJS.Timeout;
  pollingInterval = 60;
  loading = false;
  intervalOptions = [10, 30, 60, 120];

  constructor(private api: ApiService, private cdRef: ChangeDetectorRef, private subheader: SubheaderService) {
    this.subheader.setTitle('Monitor Agents');
    this.subheader.setDescription('View live data being sent from the host agents');
  }

  ngOnInit(): void {
    this.environmentSubscription = this.api.getEnvironment('a3f73ee0-da71-4aa6-9280-18ad1a1a8d16').subscribe((result) => {
      // console.log('env subscription');
      this.environment = resolveStatuses((result.data as any).environment) as Environment;
      this.loaded = true;
      this.cdRef.detectChanges();
      this.initAgentStatusPolling();
    });
  }

  initAgentStatusPolling(): void {
    // Prevent us from refetching the environment config every time
    this.environmentSubscription.unsubscribe();
    // Go ahead and query the statuses for the first time
    this.fetchAgentStatuses();
    // Set up the query to be polled every interval
    this.agentPollingInterval = setInterval(() => this.fetchAgentStatuses(), this.pollingInterval * 1000);
  }

  fetchAgentStatuses(): void {
    this.loading = true;
    this.cdRef.detectChanges();
    this.api.getAgentStatuses(this.environment.id).then((result: AgentStatusQueryResult) => {
      this.loading = false;
      this.environment = updateAgentStatuses(this.environment, result);
      this.cdRef.detectChanges();
    });
  }
  // onBranchSelect(changeEvent: MatSelectChange) {
  onIntervalChange(changeEvent: MatSelectChange): void {
    // Update the interval based on select's value
    this.pollingInterval = changeEvent.value;
    // Stop the old polling
    clearInterval(this.agentPollingInterval);
    // Set up polling again with new interval
    this.agentPollingInterval = setInterval(() => this.fetchAgentStatuses(), this.pollingInterval * 1000);
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }
}
