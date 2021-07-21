import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { UserModel } from '../../_models/user.model';
import { environment } from '../../../../../environments/environment';
import { AuthModel } from '../../_models/auth.model';
import { LaForgeAuthUser, LaForgeGetCurrentUserGQL } from '@graphql';
import { map } from 'rxjs/operators';

const API_USERS_URL = `${environment.apiUrl}/users`;

@Injectable({
  providedIn: 'root'
})
export class AuthHTTPService {
  constructor(private http: HttpClient, private getCurrentUser: LaForgeGetCurrentUserGQL) {}

  logout(): Observable<boolean> {
    return this.http
      .get(`${environment.authBaseUrl}/logout`, {
        withCredentials: true,
        observe: 'response'
      })
      .pipe(
        map((response) => {
          return response.status === 200;
        })
      );
  }

  // public methods
  login(email: string, password: string): Observable<any> {
    return this.http.post<AuthModel>(
      API_USERS_URL,
      { email, password },
      {
        withCredentials: true
      }
    );
  }

  localLogin(email: string, password: string): Observable<LaForgeAuthUser> {
    return this.http.post<LaForgeAuthUser>(
      `${environment.authBaseUrl}/local/login`,
      {
        username: email,
        password
      },
      {
        withCredentials: true
      }
    );
  }

  // CREATE =>  POST: add a new user to the server
  createUser(user: UserModel): Observable<UserModel> {
    return this.http.post<UserModel>(API_USERS_URL, user);
  }

  // Your server should check email => If email exists send link to the user and return true | If email doesn't exist return false
  forgotPassword(email: string): Observable<boolean> {
    return this.http.post<boolean>(`${API_USERS_URL}/forgot-password`, {
      email
    });
  }

  getUserByToken(token): Observable<UserModel> {
    const httpHeaders = new HttpHeaders({
      Authorization: `Bearer ${token}`
    });
    return this.http.get<UserModel>(`${API_USERS_URL}`, {
      headers: httpHeaders
    });
  }

  getCurrentUserFromContext(): Observable<LaForgeAuthUser> {
    return this.getCurrentUser.fetch().pipe(
      map((result) => {
        return result.data.currentUser;
      })
    );
  }
}
