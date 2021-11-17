import { Injectable } from '@angular/core';

export type Theme = 'light' | 'dark';

@Injectable({
  providedIn: 'root'
})
export class ThemeService {
  private currentTheme: Theme;

  constructor() {
    const cachedTheme = localStorage.getItem('laforge-theme') as Theme;
    this.currentTheme = cachedTheme || 'dark';
    this.initTheme();
  }

  public getTheme(): Theme {
    return this.currentTheme;
  }

  public setTheme(theme: Theme): void {
    this.currentTheme = theme;
    this.initTheme();
  }

  private initTheme(): void {
    localStorage.setItem('laforge-theme', this.currentTheme);

    document.body.classList.remove('theme-light');
    document.body.classList.remove('theme-dark');

    document.body.classList.add(`theme-${this.currentTheme}`);
  }
}
