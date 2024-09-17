import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SzenarioListComponent } from './szenario-list.component';

describe('SzenarioListComponent', () => {
  let component: SzenarioListComponent;
  let fixture: ComponentFixture<SzenarioListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SzenarioListComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SzenarioListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
