export class AuthModel {
  accessToken: string;
  refreshToken: string;
  expiresIn: Date;

  setAuth(auth: any) {
    this.accessToken = auth.accessToken;
    this.refreshToken = auth.refreshToken;
    this.expiresIn = auth.expiresIn;
  }
}
