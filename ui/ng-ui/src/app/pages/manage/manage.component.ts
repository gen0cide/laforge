import { Component, OnInit } from '@angular/core';
import { ProvisionedNetwork } from 'src/app/models/network.model';

import { corp_network_provisioned } from '../../../data/corp';

@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ManageComponent implements OnInit {
  corpNetwork: ProvisionedNetwork = corp_network_provisioned;

  constructor() {}

  ngOnInit(): void {}
}
