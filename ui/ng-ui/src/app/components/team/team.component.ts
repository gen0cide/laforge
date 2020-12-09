import { Component, Input, OnInit } from '@angular/core';
import { ProvisionStatus, Team } from 'src/app/models/common.model';
import { ProvisionedNetwork } from 'src/app/models/network.model';

@Component({
  selector: 'app-team',
  templateUrl: './team.component.html',
  styleUrls: ['./team.component.scss']
})
export class TeamComponent implements OnInit {
  @Input() title: string;
  @Input() networks: ProvisionedNetwork[];

  constructor() {}

  ngOnInit(): void {}

  getStatus(): ProvisionStatus {
    let status: ProvisionStatus = ProvisionStatus.ProvStatusComplete;
    for (const network of this.networks) {
      switch (network.status.state) {
        case ProvisionStatus.ProvStatusFailed:
          status = ProvisionStatus.ProvStatusFailed;
          break;
        case ProvisionStatus.ProvStatusInProgress:
          if (status === ProvisionStatus.ProvStatusComplete)
            status = ProvisionStatus.ProvStatusInProgress;
          break;
        default:
          break;
      }
      for (const host of network.provisionedHosts) {
        switch (host.status.state) {
          case ProvisionStatus.ProvStatusFailed:
            status = ProvisionStatus.ProvStatusFailed;
            break;
          case ProvisionStatus.ProvStatusInProgress:
            if (status === ProvisionStatus.ProvStatusComplete)
              status = ProvisionStatus.ProvStatusInProgress;
            break;
          default:
            break;
        }
      }
    }
    return status;
  }

  getStatusColor(): string {
    switch (this.getStatus()) {
      case ProvisionStatus.ProvStatusComplete:
        return 'success';
      case ProvisionStatus.ProvStatusInProgress:
        return 'info';
      case ProvisionStatus.ProvStatusFailed:
        return 'danger';
      default:
        return 'dark';
    }
  }
}
