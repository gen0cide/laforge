import { Component, OnInit } from '@angular/core';
import { Environment } from 'src/app/models/environment.model';
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
  environment: Environment = bradley;

  constructor() {}

  ngOnInit(): void {}
}
