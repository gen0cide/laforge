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

  constructor(private rebuild: RebuildService) {
    if (!this.style) this.style = 'compact';
    if (!this.mode) this.mode = 'manage';
    if (!this.selectable) this.selectable = false;
  }

  ngOnInit(): void {}

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
