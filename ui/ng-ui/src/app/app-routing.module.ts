import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AuthGuard } from './modules/auth/_services/auth.guard';
import { ErrorsComponent } from './modules/errors/errors.component';

export const routes: Routes = [
  {
    path: 'auth',
    loadChildren: () => import('./modules/auth/auth.module').then((m) => m.AuthModule)
  },
  {
    path: 'errors',
    component: ErrorsComponent
  },
  {
    path: '',
    canActivate: [AuthGuard],
    loadChildren: () => import('./pages/layout.module').then((m) => m.LayoutModule)
  },
  { path: '**', redirectTo: 'errors', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
