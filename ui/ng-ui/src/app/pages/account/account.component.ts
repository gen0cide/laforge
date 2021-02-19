import { Component, OnInit } from '@angular/core';
import { User } from 'src/app/models/user.model';
import { UserService } from 'src/app/services/user/user.service';
import { SubheaderService } from 'src/app/_metronic/partials/layout/subheader/_services/subheader.service';

@Component({
  selector: 'app-account',
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.scss']
})
export class AccountComponent implements OnInit {
  public user: User;

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
