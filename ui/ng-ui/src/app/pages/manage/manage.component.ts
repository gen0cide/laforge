import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
// import { ProvisionedNetwork } from 'src/app/models/network.model';

// import { corp_network_provisioned } from '../../../data/corp';
// import { bradley } from 'src/data/sample-config';
@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ManageComponent implements OnInit {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  environment: Environment = null;
  loaded = false;
  displayedColumns: string[] = ['TeamCount', 'AdminCIDRs', 'ExposedVDIPorts', 'maintainer'];

  constructor(private api: ApiService, private cdRef: ChangeDetectorRef) {}

  ngOnInit(): void {
    this.api.getEnvironment('a3f73ee0-da71-4aa6-9280-18ad1a1a8d16').subscribe((result) => {
      console.log('subsribe call');
      this.environment = resolveStatuses((result.data as any).environment) as Environment;
      this.loaded = true;
      this.cdRef.detectChanges();
    });
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }
}
