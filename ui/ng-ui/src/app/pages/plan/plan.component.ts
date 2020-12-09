import { Component, NgModule, OnInit } from '@angular/core';
import { Environment } from 'src/app/models/environment.model';

// import { Step } from '../.../../../models/plan.model';
import { PlanService } from '../../plan.service';
import {chike} from "../../../data/sample-config";

// interface EnvConfig {
//   text : string
// }

@Component({
  selector: 'app-plan',
  templateUrl: './plan.component.html',
  styleUrls: ['./plan.component.scss']
})
export class PlanComponent implements OnInit {

  environment: Environment = chike;
  constructor() {}
  ngOnInit(): void {}
}
