import { NgModule } from '@angular/core';

import { FromBytesPipe } from './bytes.pipe';

@NgModule({
  imports: [],
  declarations: [FromBytesPipe],
  exports: [FromBytesPipe]
})
export class LaforgePipesModule {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  static forRoot(): any {
    return {
      ngModule: LaforgePipesModule,
      providers: []
    };
  }
}
