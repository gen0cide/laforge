import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

// import { Step } from './models/plan.model';

@Injectable({
  providedIn: 'root'
})
export class PlanService {
  // planSteps : Step[] = [];
  planText: string;

  constructor(private http: HttpClient) {}

  clearPlan() {
    // this.planSteps = [];
    this.planText = '';
  }

  getPlan(gitBranch: string): void {
    // this.http.get('./data/tempPlan.txt', { responseType: "text" }).subscribe(rawSteps => {
    //   this.planSteps = rawSteps.split("\n").map(step => {
    //     let stepParts = step.split(' ');
    //     return {
    //       method: stepParts[0],
    //       path: stepParts[0] === "[LAFORGE:cli]" ? "" : stepParts[3]
    //     }
    //   });
    //   this.planText = rawSteps;
    //   console.log("Got plan successfully...")
    // });
  }
}
