import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { SomService } from '../../services/som.service';
import { overview } from '../../som/szenario';

@Component({
  selector: 'app-szenario-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './szenario-list.component.html',
  styleUrl: './szenario-list.component.css'
})
export class SzenarioListComponent {

  json: overview = {}

  constructor(private som: SomService) { }

  ngOnInit(): void {
    this.som.getSzenarions().subscribe(v => {
      console.log(v);
      
      this.json = v;
    })
  }
}
