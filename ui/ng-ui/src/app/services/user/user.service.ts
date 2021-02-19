import { Injectable } from '@angular/core';
import { PromiseType } from 'protractor/built/plugins';
import { User } from 'src/app/models/user.model';
import { ApiService } from '../api/api.service';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private user: User;

  constructor(private api: ApiService) {
    // TODO: refresh user session with refreshToken
  }

  public login(emailOrUsername: string, password: string): Promise<User> {
    // TODO: Perform user auth
    return Promise.reject({ message: 'Feature not implemented yet' });
  }

  public signout() {
    // TODO: Sign out user by clearing the cookies or something
    this.user = null;
    window.location.reload(); // reload window to reset them to auth screen
  }

  public me(): User {
    return this.user;
  }

  public updateUser(user: User) {
    // TODO: update the user in the DB with new details
  }

  public createUser(username: string, email: string, password: string) {
    // TODO: create the user in the DB
  }

  public getUsers(): Promise<User[]> {
    // TODO: Pull a list of the users from the DB
    return Promise.reject({ message: 'Feature not implemented yet' });
  }

  public deleteUser(): Promise<{ success: boolean }> {
    // TODO: Delete the user in the DB (possivle just archive them and not "delete" any data?)
    return Promise.reject({ message: 'Feature not implemented yet' });
  }
}
