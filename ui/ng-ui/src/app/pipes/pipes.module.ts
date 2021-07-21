import { NgModule } from '@angular/core';

import { FromBytesPipe } from './bytes.pipe';
import { DateAgoPipe } from './date-ago';
import { SortByPipe } from './sort-by';

@NgModule({
  imports: [],
  declarations: [FromBytesPipe, SortByPipe, DateAgoPipe],
  exports: [FromBytesPipe, SortByPipe, DateAgoPipe]
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
