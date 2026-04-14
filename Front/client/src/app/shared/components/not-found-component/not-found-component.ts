import { Location } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-not-found-component',
  imports: [],
  templateUrl: './not-found-component.html',
  styleUrl: './not-found-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NotFoundComponent {
  route = inject(Router);
  location = inject(Location);
  
  goBack(){
    this.location.back();
  }
}
