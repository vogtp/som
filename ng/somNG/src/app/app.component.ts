import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SzenarioListComponent } from './components/szenario-list/szenario-list.component';
import { SomService } from './services/som.service';
import {MatButtonModule, MatCheckboxModule, } from '@angular/material';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    HttpClientModule,
    SzenarioListComponent,
    BrowserAnimationsModule,
    
  ],
  providers: [SomService],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'SOM NG';
}
