import { ChangeDetectionStrategy, Component,  input, output } from '@angular/core';
import { AuthStateService } from '../../../../../auth/auth-state.service';

@Component({
  selector: 'app-overview-component',
  imports: [],
  templateUrl: './overview-component.html',
  styleUrl: './overview-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OverviewComponent {
  // userOverview = input.required();
  numberOfFiles = input.required<number>();
  typeFilesTop = input.required<{ type: string; count: number }[]>();
  currentUser = input.required();
  goToFiles = output();

  onGoToFiles() {
    this.goToFiles.emit();
  }

}
