import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NetworkComponent } from './network/network.component';
import { HostComponent } from './host/host.component';
import { TeamComponent } from './team/team.component';

@NgModule({
  declarations: [NetworkComponent, HostComponent, TeamComponent],
  imports: [CommonModule],
  exports: [NetworkComponent, HostComponent, TeamComponent]
})
export class ViewComponentsModule {}
