import { Injectable, OnDestroy } from '@angular/core';
import { Observable, BehaviorSubject, of, Subscription } from 'rxjs';
import { map, catchError, switchMap, finalize } from 'rxjs/operators';
import { UserModel } from '../_models/user.model';
import { AuthModel } from '../_models/auth.model';
import { AuthHTTPService } from './auth-http';
import { environment } from 'src/environments/environment';
import { ActivatedRoute, Router } from '@angular/router';
import { LaForgeAuthUser } from '@graphql';

@Injectable({
  providedIn: 'root'
})
export class AuthService implements OnDestroy {
  // private fields
  private unsubscribe: Subscription[] = []; // Read more: => https://brianflove.com/2016/12/11/anguar-2-unsubscribe-observables/
  private isLoadingSubject: BehaviorSubject<boolean>;
  private authLocalStorageToken = `${environment.appVersion}-${environment.USERDATA_KEY}`;

  // public fields
  // currentUser$: Observable<UserModel>;
  isLoading: Observable<boolean>;
  currentUser: Observable<LaForgeAuthUser>;
  currentUserSubject: BehaviorSubject<LaForgeAuthUser>;

  get currentUserValue(): LaForgeAuthUser {
    return this.currentUserSubject.value;
  }

  constructor(private authHttpService: AuthHTTPService, private router: Router, private route: ActivatedRoute) {
    this.isLoadingSubject = new BehaviorSubject<boolean>(false);
    this.currentUserSubject = new BehaviorSubject<LaForgeAuthUser>(undefined);
    this.currentUser = this.currentUserSubject.asObservable();
    this.isLoading = this.isLoadingSubject.asObservable();
    const subscr = this.getCurrentUserFromContext().subscribe();
    this.unsubscribe.push(subscr);
  }

  // public methods
  localLogin(email: string, password: string): Observable<LaForgeAuthUser> {
    this.isLoadingSubject.next(true);
    return this.authHttpService.localLogin(email, password).pipe(
      map((auth: LaForgeAuthUser) => {
        if (auth.id) {
          this.currentUserSubject = new BehaviorSubject<LaForgeAuthUser>(auth);
          return true;
        } else return false;
      }),
      catchError((err) => {
        console.error('err', err);
        return of(undefined);
      }),
      finalize(() => {
        this.isLoadingSubject.next(false);
      })
    );
  }

  logout(): Observable<boolean> {
    return this.authHttpService.logout().pipe(
      map((success) => {
        return success;
      })
    );
  }

  getCurrentUserFromContext(): Observable<LaForgeAuthUser> {
    this.isLoadingSubject.next(true);
    return this.authHttpService.getCurrentUserFromContext().pipe(
      map((user: LaForgeAuthUser) => {
        if (user) {
          this.currentUserSubject = new BehaviorSubject<LaForgeAuthUser>(user);
        } else {
          this.logout();
        }
        return user;
      }),
      catchError((err) => {
        console.warn('err', err);
        this.router.navigate(['auth', 'login']);
        return of(undefined);
      }),
      finalize(() => this.isLoadingSubject.next(false))
    );
  }

  // need create new user then login
  registration(user: UserModel): Observable<any> {
    this.isLoadingSubject.next(true);
    return this.authHttpService.createUser(user).pipe(
      map(() => {
        this.isLoadingSubject.next(false);
      }),
      switchMap(() => this.localLogin(user.email, user.password)),
      catchError((err) => {
        console.error('err', err);
        return of(undefined);
      }),
      finalize(() => this.isLoadingSubject.next(false))
    );
  }

  forgotPassword(email: string): Observable<boolean> {
    this.isLoadingSubject.next(true);
    return this.authHttpService.forgotPassword(email).pipe(finalize(() => this.isLoadingSubject.next(false)));
  }
  // private methods
  private setAuthFromLocalStorage(auth: AuthModel): boolean {
    // store auth accessToken/refreshToken/epiresIn in local storage to keep user logged in between page refreshes
    if (auth && auth.accessToken) {
      localStorage.setItem(this.authLocalStorageToken, JSON.stringify(auth));
      return true;
    }
    return false;
  }

  private getAuthFromLocalStorage(): AuthModel {
    try {
      const authData = JSON.parse(localStorage.getItem(this.authLocalStorageToken));
      return authData;
    } catch (error) {
      console.error(error);
      return undefined;
    }
  }

  ngOnDestroy() {
    this.unsubscribe.forEach((sb) => sb.unsubscribe());
  }
}
