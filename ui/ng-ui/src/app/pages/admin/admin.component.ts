import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { LaForgeAuthUser, LaForgeGetUserListQuery } from '@graphql';
import { ApiService } from '@services/api/api.service';
import { EnvironmentService } from '@services/environment/environment.service';
import { BehaviorSubject } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { AuthService } from 'src/app/modules/auth/_services/auth.service';

import { EditUserModalComponent } from '@components/edit-user-modal/edit-user-modal.component';

@Component({
  selector: 'app-dashboard',
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.scss']
})
export class AdminComponent implements OnInit {
  public loading: BehaviorSubject<boolean>;
  public errorMessage: BehaviorSubject<string>;
  public userList: BehaviorSubject<LaForgeGetUserListQuery['getUserList']>;
  public userListCols: string[] = ['name', 'username', 'email', 'provider', 'role', 'controls'];

  constructor(
    public dialog: MatDialog,
    private cdRef: ChangeDetectorRef,
    private subheader: SubheaderService,
    public envService: EnvironmentService,
    private api: ApiService,
    private auth: AuthService,
    private snackbar: MatSnackBar,
    private router: Router
  ) {
    this.subheader.setTitle('Manage Users');
    this.subheader.setDescription('Add, delete, or modify users');
    this.subheader.setShowEnvDropdown(false);

    this.loading = new BehaviorSubject<boolean>(false);
    this.errorMessage = new BehaviorSubject<string>('');
    this.userList = new BehaviorSubject<LaForgeGetUserListQuery['getUserList']>([]);
  }

  ngOnInit(): void {
    this.auth.getCurrentUserFromContext().subscribe(
      (user) => {
        if (user.role === 'ADMIN') {
          this.refreshUserList();
        } else {
          this.router.navigate(['dashboard']);
        }
      },
      () => this.router.navigate(['dashboard'])
    );

    // this.dialog.afterAllClosed.subscribe(() => window.location.reload());
  }

  refreshUserList() {
    this.loading.next(true);
    this.errorMessage.next('');
    this.api
      .getAllUsers()
      .then(
        (users) => {
          this.userList.next(users);
        },
        (err) => {
          this.userList.error(err);
          this.errorMessage.next(err);
          this.snackbar.open(err, 'Okay', {
            panelClass: ['bg-danger', 'text-white']
          });
        }
      )
      .finally(() => this.loading.next(false));
  }

  editUser(user: LaForgeAuthUser) {
    MatDialogConfig;
    this.dialog
      .open(EditUserModalComponent, {
        data: {
          user: user
        },
        width: '50%',
        height: '75%',
        autoFocus: true
      })
      .afterClosed()
      .subscribe(() => this.refreshUserList());
  }

  createUser() {
    this.dialog
      .open(EditUserModalComponent, {
        data: {
          user: null
        },
        width: '50%',
        height: '75%',
        autoFocus: true
      })
      .afterClosed()
      .subscribe(() => this.refreshUserList());
  }
}
