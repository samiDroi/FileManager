import { ChangeDetectionStrategy, Component, output, input } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { debounceTime, distinctUntilChanged, Subject } from 'rxjs';

@Component({
  selector: 'app-ui-search-bar',
  imports: [FormsModule],
  templateUrl: './ui-search-bar-component.html',
  styleUrl: './ui-search-bar-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UiSearchBarComponent {
  placeholder = input<string>('Buscar...');
  debounceMs = input<number>(300);

  search = output<string>();

  searchTerm = '';
  private searchSubject = new Subject<string>();

  constructor() {
    this.searchSubject.pipe(
      // debounceTime(this.debounceMs()),
      distinctUntilChanged()
    ).subscribe(term => {
      this.search.emit(term);
    });
  }

  onInput(value: string): void {
      this.searchTerm = value;
  this.searchSubject.next(value);
  }

  onClear(): void {
    this.searchTerm = '';
    this.search.emit('');
  }
}
