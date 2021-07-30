import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Router } from '@angular/router';

import { LaForgeDeleteBuildGQL } from '../../../generated/graphql';

@Component({
  selector: 'app-delete-build-modal',
  templateUrl: './delete-build-modal.component.html',
  styleUrls: ['./delete-build-modal.component.scss']
})
export class DeleteBuildModalComponent {
  deleteConfirmed = false;
  deleteLoading = false;

  constructor(
    public dialogRef: MatDialogRef<DeleteBuildModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { buildName: string; buildId: string },
    private deleteBuild: LaForgeDeleteBuildGQL,
    private router: Router
  ) {}

  buildNameChange(value: string) {
    if (value === this.data.buildName) {
      this.deleteConfirmed = true;
    } else {
      this.deleteConfirmed = false;
    }
  }

  onClose(): void {
    this.dialogRef.close();
  }

  triggerDelete(): void {
    if (!this.data.buildId) return;
    this.deleteLoading = true;
    this.deleteBuild
      .mutate({
        buildId: this.data.buildId
      })
      .toPromise()
      .then(({ data, errors }) => {
        if (errors) {
          return console.error(errors);
        } else if (data.deleteBuild) {
          this.router.navigate(['plan']);
          return this.onClose();
        }
        console.error('delete build failed');
      }, console.error)
      .finally(() => (this.deleteLoading = false));
  }
}
