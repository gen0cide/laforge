import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DeleteBuildModalComponent } from './delete-build-modal.component';

describe('DeleteBuildModalComponent', () => {
  let component: DeleteBuildModalComponent;
  let fixture: ComponentFixture<DeleteBuildModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [DeleteBuildModalComponent]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DeleteBuildModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
