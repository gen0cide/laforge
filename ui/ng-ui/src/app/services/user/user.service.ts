import { Injectable } from '@angular/core';
import { LaForgeAuthUser } from '@graphql';

import { Observable } from 'rxjs';

import { AuthService } from '../../modules/auth/_services/auth.service';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private user: LaForgeAuthUser;

  constructor(private authService: AuthService) {
    // TODO: refresh user session with refreshToken
  }

  public login(email: string, password: string): Observable<LaForgeAuthUser> {
    // TODO: Perform user auth
    // return Promise.reject({ message: 'Feature not implemented yet' });
    return this.authService.localLogin(email, password);
    // return this.api.fakeLogin().then(
    //   (user) => {
    //     this.user = user;
    //     return true;
    //   },
    //   () => {
    //     return false;
    //   }
    // );
  }

  public signout() {
    // TODO: Sign out user by clearing the cookies or something
    this.user = null;
    window.location.reload(); // reload window to reset them to auth screen
  }

  public me(): LaForgeAuthUser {
    return this.user;
  }

  public updateUser(user: LaForgeAuthUser) {
    // TODO: update the user in the DB with new details
  }

  public createUser(username: string, email: string, password: string) {
    // TODO: create the user in the DB
  }

  public getUsers(): Promise<LaForgeAuthUser[]> {
    // TODO: Pull a list of the users from the DB
    return Promise.reject({ message: 'Feature not implemented yet' });
  }

  public deleteUser(): Promise<{ success: boolean }> {
    // TODO: Delete the user in the DB (possivle just archive them and not "delete" any data?)
    return Promise.reject({ message: 'Feature not implemented yet' });
  }
}
