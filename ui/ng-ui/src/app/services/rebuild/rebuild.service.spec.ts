import { TestBed } from '@angular/core/testing';

import { RebuildService } from './rebuild.service';

describe('RebuildService', () => {
  let service: RebuildService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(RebuildService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
