import { Injectable } from '@angular/core';
import { ProvisionedHost } from 'src/app/models/host.model';
import { ProvisionedNetwork } from 'src/app/models/network.model';

@Injectable({
  providedIn: 'root'
})
export class RebuildService {
  networksToRebuild: ProvisionedNetwork[];
  hostsToRebuild: ProvisionedHost[];

  constructor() {
    this.networksToRebuild = [];
    this.hostsToRebuild = [];
  }

  addHost = (host: ProvisionedHost): void => {
    if (this.hostsToRebuild.filter((h) => h.id === host.id).length === 0) this.hostsToRebuild.push(host);
  };

  removeHost = (host: ProvisionedHost): void => {
    if (this.hostsToRebuild.indexOf(host) >= 0) this.hostsToRebuild.splice(this.hostsToRebuild.indexOf(host), 1);
  };

  addNetwork = (network: ProvisionedNetwork): void => {
    if (this.networksToRebuild.filter((n) => n.id === network.id).length === 0) {
      this.networksToRebuild.push(network);
      network.provisionedHosts.forEach((host) => {
        if (this.hostsToRebuild.filter((h) => h.id === host.id).length > 0) this.removeHost(host);
      });
    }
  };

  removeNetwork = (network: ProvisionedNetwork): void => {
    if (this.networksToRebuild.indexOf(network) >= 0) this.networksToRebuild.splice(this.networksToRebuild.indexOf(network), 1);
  };
}
