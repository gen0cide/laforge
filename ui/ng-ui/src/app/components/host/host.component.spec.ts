import { Renderer2, ChangeDetectorRef, Type } from '@angular/core';
import {
  async,
  ComponentFixture,
  fakeAsync,
  TestBed,
  tick
} from '@angular/core/testing';

import { HostComponent } from './host.component';

import { coins_heads_01_provisioned } from 'src/data/corp';

describe('HostComponent', () => {
  let component: HostComponent;
  let fixture: ComponentFixture<HostComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [HostComponent],
      providers: [Renderer2, ChangeDetectorRef]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HostComponent);
    component = fixture.componentInstance;
    component.provisionedHost = coins_heads_01_provisioned;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('#clicked() should open dropdown', fakeAsync(() => {
    spyOn(component, 'toggleOptions');
    const container = fixture.debugElement.nativeElement.querySelector(
      '.host-info'
    );
    const details = fixture.debugElement.nativeElement.querySelector(
      '.host-details'
    );
    expect(component.optionsToggled).toBeFalsy('details hidden at first');
    expect(details.classList.contains('expanded')).toBeFalsy(
      ".host-details doesn't have expanded class"
    );
    // component.toggleOptions();
    console.log(container);
    container.click();
    tick();
    expect(component.toggleOptions).toHaveBeenCalled();
    console.log(
      fixture.debugElement.nativeElement.querySelector('.host-details')
        .classList
    );
    expect(details.classList.contains('expanded')).toBeTruthy(
      '.host-details has expanded class'
    );
    container.click();
    tick();
    expect(component.optionsToggled).toBeFalse();
  }));
});
