
import { NgTemplateOutlet } from '@angular/common';
import { ChangeDetectionStrategy, Component, HostListener, input } from '@angular/core';
import { TemplateRef } from '@angular/core';

@Component({
  selector: 'app-ui-table',
  imports: [NgTemplateOutlet],
  templateUrl: './ui-table-component.html',
  styleUrl: './ui-table-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UiTableComponent {
  data = input<any[]>([]);
  //template for each row
  rowTemplate = input.required<TemplateRef<any>>();
  colSpan = input<number>(4);

  menuTemplate = input<TemplateRef<any>>();
  closeMenuFn = () => { this.activeMenuIndex = null; }
  //row from ctable that triggers the action (menu)
  activeMenuIndex: number | null = null;

  toggleMenu(index: number, event: Event): void {
    event.stopPropagation();
    this.activeMenuIndex = this.activeMenuIndex === index ? null : index;
  }

  @HostListener('document:click')
  closeMenu(): void {
    this.activeMenuIndex = null;
  }
}