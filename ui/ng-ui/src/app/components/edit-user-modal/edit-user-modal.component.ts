import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { LaForgeAuthUser, LaForgeRoleLevel, LaForgeProviderType } from '@graphql';
import { ApiService } from '@services/api/api.service';
import { BehaviorSubject } from 'rxjs';

@Component({
  selector: 'app-edit-user-modal',
  templateUrl: './edit-user-modal.component.html',
  styleUrls: ['./edit-user-modal.component.scss']
})
export class EditUserModalComponent implements OnInit {
  username = new FormControl('', Validators.required);
  password = new FormControl('', Validators.required);
  firstName = new FormControl('');
  lastName = new FormControl('');
  email = new FormControl('', [Validators.required, Validators.email]);
  phone = new FormControl('', Validators.pattern(/[0-9\+\-\ ]/));
  company = new FormControl('');
  occupation = new FormControl('');
  role = new FormControl('', Validators.required);
  roleList: string[] = [LaForgeRoleLevel.Admin, LaForgeRoleLevel.User];
  provider = new FormControl('', Validators.required);
  providerList: string[] = [LaForgeProviderType.Local, LaForgeProviderType.Github, LaForgeProviderType.Openid];
  errorMessage: BehaviorSubject<string>;

  constructor(
    public dialogRef: MatDialogRef<EditUserModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { user?: LaForgeAuthUser },
    private api: ApiService
  ) {}

  ngOnInit(): void {
    console.log(this.data.user);
    this.username.setValue(this.data.user?.username ?? '');
    this.firstName.setValue(this.data.user?.first_name ?? '');
    this.lastName.setValue(this.data.user?.last_name ?? '');
    this.email.setValue(this.data.user?.email ?? '');
    this.phone.setValue(this.data.user?.phone ?? '');
    this.company.setValue(this.data.user?.company ?? '');
    this.occupation.setValue(this.data.user?.occupation ?? '');
    this.role.setValue(this.data.user?.role ?? '');
    this.provider.setValue(this.data.user?.provider ?? 'LOCAL');

    this.errorMessage = new BehaviorSubject<string>('');
  }

  getUsernameErrorMessage(): string {
    if (this.email.hasError('required')) {
      return 'This field is required';
    }
    return '';
  }

  getPasswordErrorMessage(): string {
    if (this.email.hasError('required')) {
      return 'This field is required';
    }
    return '';
  }

  getEmailErrorMessage(): string {
    if (this.email.hasError('required')) {
      return 'This field is required';
    }
    if (this.email.hasError('email')) {
      return 'Must be a valid email';
    }
    return '';
  }

  getPhoneErrorMessage(): string {
    if (this.phone.hasError('pattern')) {
      return 'Must be a phone number (no spaces or parens)';
    }
    return '';
  }

  getRoleErrorMessage(): string {
    if (this.role.hasError('required')) {
      return 'This field is required';
    }
    return '';
  }

  getProviderErrorMessage(): string {
    if (this.provider.hasError('required')) {
      return 'This field is required';
    }
    return '';
  }

  cancel() {
    this.dialogRef.close();
  }

  submit() {
    if (this.data.user) {
      this.api
        .modifyUser(this.data.user.id, {
          firstName: this.firstName.value,
          lastName: this.lastName.value,
          email: this.email.value,
          phone: this.phone.value,
          company: this.company.value,
          occupation: this.occupation.value,
          role: this.role.value,
          provider: this.provider.value
        })
        .then(
          () => {
            this.dialogRef.close();
          },
          (err) => this.errorMessage.next(err)
        );
    } else {
      this.api
        .createUser({
          username: this.username.value,
          password: this.password.value,
          role: this.role.value,
          provider: this.provider.value
        })
        .then(
          () => {
            this.dialogRef.close();
          },
          (err) => this.errorMessage.next(err)
        );
    }
  }
}
