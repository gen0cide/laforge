import { Component, OnInit } from '@angular/core';
import { Environment } from 'src/app/models/environment.model';
import { ApiService } from 'src/app/services/api/api.service';
// import { ProvisionedNetwork } from 'src/app/models/network.model';

// import { corp_network_provisioned } from '../../../data/corp';
import { bradley } from 'src/data/sample-config';
@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ManageComponent implements OnInit {
  // corpNetwork: ProvisionedNetwork = corp_network_provisioned;
  environment: Environment;

  constructor(private api: ApiService) {}

  ngOnInit(): void {
    this.api.getEnvironment('a3f73ee0-da71-4aa6-9280-18ad1a1a8d16').then((environment: Environment) => {
      this.environment = environment;
      // console.log(environment);
    });
  }
}
