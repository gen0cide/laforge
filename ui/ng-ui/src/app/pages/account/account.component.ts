import { Component, OnInit } from '@angular/core';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';
import { UserService } from 'src/app/services/user/user.service';

import { AuthUser } from '../../models/user.model';

@Component({
  selector: 'app-account',
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.scss']
})
export class AccountComponent implements OnInit {
  public user: AuthUser;

  constructor(private subheader: SubheaderService, private userService: UserService) {
    this.subheader.setTitle('Account');
    this.subheader.setDescription('Edit your account details');
    this.user = this.userService.me();
  }

  ngOnInit(): void {}

  onSubmitChanges(event: Event) {
    event.preventDefault();
    this.userService.updateUser(this.user);
  }
}
