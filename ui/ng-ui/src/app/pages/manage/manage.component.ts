import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
import { SubheaderService } from '../../_metronic/partials/layout/subheader/_services/subheader.service';

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
  selectionMode = false;

  constructor(private api: ApiService, private cdRef: ChangeDetectorRef, private subheader: SubheaderService) {
    this.subheader.setTitle('Environment');
    this.subheader.setDescription('Manage your currently running environment');
  }

  ngOnInit(): void {
    this.api.getEnvironment('a3f73ee0-da71-4aa6-9280-18ad1a1a8d16').subscribe((result) => {
      this.environment = resolveStatuses(result.data.environment) as Environment;
      this.loaded = true;
      this.cdRef.detectChanges();
    });
  }

  rebuildEnv(): void {
    console.log('rebuilding env...');
  }

  toggleSelectionMode(): void {
    this.selectionMode = !this.selectionMode;
  }
}
