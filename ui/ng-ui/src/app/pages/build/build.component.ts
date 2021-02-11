import { Component, OnInit } from '@angular/core';
import { Environment } from 'src/app/models/environment.model';
import { EnvironmentService } from 'src/app/services/environment/environment.service';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-build',
  templateUrl: './build.component.html',
  styleUrls: ['./build.component.scss']
})
export class BuildComponent implements OnInit {
  environment: Observable<Environment>;

  constructor(private subheader: SubheaderService, public envService: EnvironmentService) {
    this.subheader.setTitle('Build');
    this.subheader.setDescription('Build a planned environment');

    this.environment = this.envService.getCurrentEnv().asObservable();
  }

  ngOnInit(): void {}

  envIsSelected(): boolean {
    return this.envService.getCurrentEnv().getValue() != null;
  }
}
