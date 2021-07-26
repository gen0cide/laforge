import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { EnvironmentService } from '@services/environment/environment.service';
import { Observable } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  getEnvironmentsLoading: Observable<boolean>;

  constructor(private cdRef: ChangeDetectorRef, private subheader: SubheaderService, public envService: EnvironmentService) {
    this.subheader.setTitle('Dashboard');
    this.subheader.setDescription('Overview of all environments and builds');
    this.subheader.setShowEnvDropdown(false);

    this.getEnvironmentsLoading = this.envService.envIsLoading.asObservable();
  }

  ngOnInit(): void {
    // this.envService.getEnvironments().subscribe(() => {
    //   this.cdRef.markForCheck();
    // });
  }
}
