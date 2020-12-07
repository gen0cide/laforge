import { Component, NgModule, OnInit } from '@angular/core';

// import { Step } from '../.../../../models/plan.model';
import { PlanService } from '../../plan.service';

// interface EnvConfig {
//   text : string
// }

@Component({
  selector: 'app-plan',
  templateUrl: './plan.component.html',
  styleUrls: ['./plan.component.scss']
})
export class PlanComponent implements OnInit {
  // currentConfig : EnvConfig = {
  //   text: 'This is a test config\
  //   More config here\
  //   Last line of the config'
  // }

  // planSteps : Step[] = []
  // planText : String

  constructor(
    public planService: PlanService // TODO: Adjust scope so we can access values while this is private
  ) {}

  ngOnInit(): void {
    // this.planService.planSteps
    // this.planService.getStepsObserver().subscribe(
    //   (steps : Step[]) => {
    //     this.planSteps = steps;
    //     this.planText = steps.map(step => step.method + " | " + step.path).reduce((prev, curr) => prev + curr + "\n", "");
    //     console.log("plan set")
    //   },
    //   (error) => {
    //     console.log(error)
    //   }
    // )
  }
}
