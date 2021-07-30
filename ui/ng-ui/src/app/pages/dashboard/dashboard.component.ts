import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LaForgeGetEnvironmentsQuery } from '@graphql';
import { ApiService } from '@services/api/api.service';
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
  // Validate the gitUrl input is a github ssh url
  gitUrl = new FormControl('', [Validators.required, Validators.pattern('(.*?)@(.*?):(?:(.*?)/)?(.*?/.*?)')]);
  repoName = new FormControl('', Validators.required);
  branchName = new FormControl('', Validators.required);
  envFilePath = new FormControl('', Validators.required);
  environmentsCols: string[] = ['name', 'competition_id', 'build_count', 'actions'];
  environments: Observable<LaForgeGetEnvironmentsQuery['environments']>;

  constructor(
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    public envService: EnvironmentService,
    private api: ApiService,
    private snackbar: MatSnackBar
  ) {
    this.subheader.setTitle('Dashboard');
    this.subheader.setDescription('Overview of all environments and builds');
    this.subheader.setShowEnvDropdown(false);

    this.getEnvironmentsLoading = this.envService.envIsLoading.asObservable();
    this.environments = this.envService.getEnvironments().asObservable();
  }

  ngOnInit(): void {
    // this.envService.getEnvironments().subscribe(() => {
    //   this.cdRef.markForCheck();
    // });
  }

  getGitErrorMessage(): string {
    if (this.gitUrl.hasError('required')) {
      return 'This field is required';
    }
    if (this.gitUrl.hasError('pattern')) {
      return 'Git URL must be a SSH URL';
    }
    return '';
  }

  getRepoNameErrorMessage(): string {
    if (this.repoName.hasError('required')) {
      return 'This field is required';
    }
    return '';
  }

  getBranchNameErrorMessage(): string {
    if (this.branchName.hasError('required')) {
      return 'This field is required';
    }
    return '';
  }

  getEnvFilePathErrorMessage(): string {
    if (this.envFilePath.hasError('required')) {
      return 'This field is required';
    }
    return '';
  }

  createBuild(envId: string) {
    this.api.createBuild(envId).then(
      (build) => {
        if (build.id) {
          this.snackbar.open('Build created. Please wait for files to render.', 'Okay', {
            duration: 3000,
            panelClass: ['bg-success', 'text-white']
          });
          this.envService.initEnvironments();
          this.envService.setCurrentEnv(envId, build.id);
        }
      },
      (err) => {
        console.error(err);
        this.snackbar.open('Error while creating build. Please check logs for details.', 'Okay', {
          duration: 3000,
          panelClass: ['bg-danger', 'text-white']
        });
      }
    );
  }

  cloneEnvironmentFromGit() {
    if (this.gitUrl.errors) return;
    if (this.repoName.errors) return;
    if (this.branchName.errors) return;
    if (this.envFilePath.errors) return;
    console.log('clone env');
  }
}
