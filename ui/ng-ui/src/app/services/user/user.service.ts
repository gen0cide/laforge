import { Injectable } from '@angular/core';
import { AuthUser, User } from '@models/user.model';

import { ApiService } from '@services/api/api.service';
import { PromiseType } from 'protractor/built/plugins';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private user: AuthUser;

  constructor(private api: ApiService) {
    // TODO: refresh user session with refreshToken
  }

  public login(email: string, password: string): Promise<boolean> {
    // TODO: Perform user auth
    // return Promise.reject({ message: 'Feature not implemented yet' });
    return this.api.fakeLogin().then(
      (user) => {
        this.user = user;
        return true;
      },
      () => {
        return false;
      }
    );
  }

  public signout() {
    // TODO: Sign out user by clearing the cookies or something
    this.user = null;
    window.location.reload(); // reload window to reset them to auth screen
  }

  public me(): AuthUser {
    return this.user;
  }

  public updateUser(user: AuthUser) {
    // TODO: update the user in the DB with new details
  }

  public createUser(username: string, email: string, password: string) {
    // TODO: create the user in the DB
  }

  public getUsers(): Promise<AuthUser[]> {
    // TODO: Pull a list of the users from the DB
    return Promise.reject({ message: 'Feature not implemented yet' });
  }

  public deleteUser(): Promise<{ success: boolean }> {
    // TODO: Delete the user in the DB (possivle just archive them and not "delete" any data?)
    return Promise.reject({ message: 'Feature not implemented yet' });
  }
}
