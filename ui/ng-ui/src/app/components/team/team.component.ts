import { PortalInjector } from '@angular/cdk/portal';
import { Component, Input, OnInit } from '@angular/core';
import { ProvisionStatus, Team } from 'src/app/models/common.model';
import { RebuildService } from '../../services/rebuild/rebuild.service';

@Component({
  selector: 'app-team',
  templateUrl: './team.component.html',
  styleUrls: ['./team.component.scss']
})
export class TeamComponent implements OnInit {
  @Input() title: string;
  @Input() team: Team;
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() mode: 'plan' | 'build' | 'manage';
  isSelected = false;

  constructor(private rebuild: RebuildService) {}

  ngOnInit(): void {
    if (!this.mode) this.mode = 'manage';
    if (!this.style) {
      switch (this.mode) {
        case 'plan':
          this.style = 'expanded';
          break;
        case 'build':
          // this.style = 'expanded'
          if (this.team) {
            switch (this.getStatus()) {
              case ProvisionStatus.ProvStatusComplete:
                this.style = 'collapsed';
                break;
              default:
                this.style = 'expanded';
                break;
            }
          }
          break;
        case 'manage':
          this.style = 'compact';
          break;
        default:
          this.style = 'compact';
          break;
      }
    }
    if (!this.selectable) this.selectable = false;
  }

  getStatus(): ProvisionStatus {
    // let status: ProvisionStatus = ProvisionStatus.ProvStatusUndefined;
    let numWithAgentData = 0;
    let totalAgents = 0;
    for (const network of this.team.provisionedNetworks) {
      for (const host of network.provisionedHosts) {
        totalAgents++;
        if (host.heartbeat) numWithAgentData++;
      }
    }
    if (numWithAgentData === totalAgents) return ProvisionStatus.ProvStatusComplete;
    else if (numWithAgentData === 0) return ProvisionStatus.ProvStatusFailed;
    else return ProvisionStatus.ProvStatusInProgress;
  }

  getStatusColor(): string {
    switch (this.getStatus()) {
      case ProvisionStatus.ProvStatusComplete:
        return 'success';
      case ProvisionStatus.ProvStatusInProgress:
        return 'warning';
      case ProvisionStatus.ProvStatusFailed:
        return 'danger';
      default:
        return 'dark';
    }
  }

  onSelect() {
    this.isSelected = !this.isSelected;
    if (this.isSelected) this.team.provisionedNetworks.forEach(this.rebuild.addNetwork);
    else this.team.provisionedNetworks.forEach(this.rebuild.removeNetwork);
  }
}
