import { Injectable } from '@angular/core';
import { EnvironmentInfo } from 'src/app/models/api.model';

@Injectable({
  providedIn: 'root'
})
export class EnvironmentService {
  private currEnvironment: EnvironmentInfo = null;
  constructor() {}

  public getCurrentEnv(): EnvironmentInfo {
    return this.currEnvironment;
  }

  public setCurrentEnv(env: EnvironmentInfo): void {
    this.currEnvironment = env;
  }
}
