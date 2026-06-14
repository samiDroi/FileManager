import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-form-component',
  imports: [],
  templateUrl: './form-component.html',
  styleUrl: './form-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FormComponent {}
