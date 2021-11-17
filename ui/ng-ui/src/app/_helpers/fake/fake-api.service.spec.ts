import { TestBed } from '@angular/core/testing';

import { FakeAPIService } from './fake-api.service';

describe('FakeAPIService', () => {
  let service: FakeAPIService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FakeAPIService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
