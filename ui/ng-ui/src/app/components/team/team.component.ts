import { PortalInjector } from '@angular/cdk/portal';
import { Component, Input, OnInit } from '@angular/core';
import { LaForgeGetBuildTreeQuery, LaForgeTeam } from '@graphql';
import { teamChildrenCompleted } from '@util';
import { ProvisionStatus } from 'src/app/models/common.model';

import { RebuildService } from '../../services/rebuild/rebuild.service';

@Component({
  selector: 'app-team',
  templateUrl: './team.component.html',
  styleUrls: ['./team.component.scss']
})
export class TeamComponent implements OnInit {
  @Input() title: string;
  @Input() team: LaForgeTeam;
  @Input() style: 'compact' | 'collapsed' | 'expanded';
  @Input() selectable: boolean;
  @Input() mode: 'plan' | 'build' | 'manage';
  isSelectedState = false;

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
              case ProvisionStatus.COMPLETE:
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
    for (const network of this.team.TeamToProvisionedNetwork) {
      for (const host of network.ProvisionedNetworkToProvisionedHost) {
        totalAgents++;
        if (host.ProvisionedHostToAgentStatus?.clientId) numWithAgentData++;
      }
    }
    if (numWithAgentData === totalAgents) return ProvisionStatus.COMPLETE;
    else if (numWithAgentData === 0) return ProvisionStatus.FAILED;
    else return ProvisionStatus.INPROGRESS;
  }

  getStatusColor(): string {
    switch (this.getStatus()) {
      case ProvisionStatus.COMPLETE:
        return 'success';
      case ProvisionStatus.INPROGRESS:
        return 'warning';
      case ProvisionStatus.FAILED:
        return 'danger';
      default:
        return 'dark';
    }
  }

  onSelect(): void {
    let success = false;
    if (!this.isSelected()) {
      success = this.rebuild.addTeam(this.team);
    } else {
      success = this.rebuild.removeTeam(this.team);
    }
    if (success) this.isSelectedState = !this.isSelectedState;
  }

  isSelected(): boolean {
    return this.rebuild.hasTeam(this.team);
  }

  shouldCollapse(): boolean {
    return teamChildrenCompleted(this.team);
  }
}
