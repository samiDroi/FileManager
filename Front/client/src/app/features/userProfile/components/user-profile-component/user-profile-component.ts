import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-user-profile-component',
  imports: [],
  templateUrl: './user-profile-component.html',
  styleUrl: './user-profile-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UserProfileComponent {}
