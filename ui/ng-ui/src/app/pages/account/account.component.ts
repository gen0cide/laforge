import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';

import { MatSnackBar } from '@angular/material/snack-bar';
import { LaForgeAuthUser, LaForgeRoleLevel } from '@graphql';
import { ApiService } from '@services/api/api.service';
import { BehaviorSubject } from 'rxjs';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';

import { AuthService } from 'src/app/modules/auth/_services/auth.service';
@Component({
  selector: 'app-account',
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.scss']
})
export class AccountComponent implements OnInit {
  public user: BehaviorSubject<LaForgeAuthUser>;
  /*
id
    username
    role
    provider
    first_name
    last_name
    email
    phone
    company
    occupation
    publicKey
  */
  username = new FormControl('', Validators.required);
  firstName = new FormControl('');
  lastName = new FormControl('');
  email = new FormControl('', Validators.email);
  phone = new FormControl('');
  company = new FormControl('');
  occupation = new FormControl('');
  publicKey = new FormControl('');

  constructor(
    private subheader: SubheaderService,
    private auth: AuthService,
    private cdRef: ChangeDetectorRef,
    private api: ApiService,
    private snackbar: MatSnackBar
  ) {
    this.subheader.setTitle('Account');
    this.subheader.setDescription('Edit your account details');
    this.user = new BehaviorSubject(null);
  }

  ngOnInit(): void {
    this.user.next(this.auth.currentUserValue);
    this.username.setValue(this.auth.currentUserValue.username);
    this.firstName.setValue(this.auth.currentUserValue.first_name);
    this.lastName.setValue(this.auth.currentUserValue.last_name);
    this.email.setValue(this.auth.currentUserValue.email);
    this.phone.setValue(this.auth.currentUserValue.phone);
    this.company.setValue(this.auth.currentUserValue.company);
    this.occupation.setValue(this.auth.currentUserValue.occupation);
    this.publicKey.setValue(this.auth.currentUserValue.publicKey);
    this.cdRef.detectChanges();
  }

  getRoleColor(): string {
    if (!this.user.getValue()) return 'link';
    switch (this.user.getValue().role) {
      case LaForgeRoleLevel.Admin:
        return 'warn';
      case LaForgeRoleLevel.User:
        return 'primary';
      default:
        return 'link';
    }
  }

  onSubmitChanges(event: Event) {
    event.preventDefault();
    this.api
      .updateAuthUser({
        firstName: this.firstName.value,
        lastName: this.lastName.value,
        email: this.email.value,
        phone: this.phone.value,
        company: this.company.value,
        occupation: this.occupation.value
      })
      .then(
        () => {
          this.auth
            .getCurrentUserFromContext()
            .toPromise()
            .then(() => {
              this.user.next(this.auth.currentUserValue);
              this.snackbar.open('Updated user successfully.', null, {
                duration: 1000,
                panelClass: ['bg-success', 'text-white']
              });
              this.cdRef.detectChanges();
            });
        },
        (err) => {
          console.error(err);
          this.snackbar.open('Error while updating user. See console for more info.', 'Okay', {
            duration: 3000,
            panelClass: ['bg-danger', 'text-white']
          });
        }
      );
  }
}
