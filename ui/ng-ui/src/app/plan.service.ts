import { Injectable } from '@angular/core';

import { HttpClient } from '@angular/common/http';

interface Step {
  method : string,
  path : string,
}

@Injectable({
  providedIn: 'root'
})
export class PlanService {
  planSteps : Step[] = [];

  constructor(private http : HttpClient) {}

  getPlan(gitBranch : string) {
    let rawSteps = this.http.get('/data/tempPlan.txt', { responseType: "text" });
    this.planSteps = rawSteps.split("\n").map(step => {
      let stepParts = step.split(' ');
      return {
        method: stepParts[0],
        path: stepParts[0] === "[LAFORGE:cli]" ? "" : stepParts[3]
      }
    })
  }
}
