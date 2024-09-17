import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { overview } from '../som/szenario';

const httpOptions = {
  headers: new HttpHeaders({
    'Accept': 'application/json',
    'Access-Control-Allow-Origin': '*'
  }),
};

@Injectable({
  providedIn: 'root'
})
export class SomService {
  private somUrl = "http://localhost:4200/som"

  constructor(private http: HttpClient) { }


  getSzenarions(): Observable<overview> {
    return this.http.get<overview>(this.somUrl, httpOptions)
  }

}
