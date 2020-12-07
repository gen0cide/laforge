import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NetworkComponent } from './network/network.component';
import { HostComponent } from './host/host.component';

@NgModule({
  declarations: [NetworkComponent, HostComponent],
  imports: [CommonModule],
  exports: [NetworkComponent]
})
export class ViewComponentsModule {}
