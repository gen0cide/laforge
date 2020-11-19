import { Component, NgModule, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';

interface Branch {
  name : string,
  hash : string
}

interface EnvConfig {
  text : string
}

@Component({
  selector: 'app-plan',
  templateUrl: './plan.component.html',
  styleUrls: ['./plan.component.scss']
})
export class PlanComponent implements OnInit {
  gitBranch = new FormControl('');
  branches : Branch[] = [
    {name: 'Bradley', hash: '98y3if'},
    {name: 'Lucas', hash: '32a7fh'},
  ]
  currentConfig : EnvConfig = {
    text: 'This is a test config\
    More config here\
    Last line of the config'
  }

  constructor() { }

  ngOnInit(): void {
  }

}
