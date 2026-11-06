import { ChangeDetectionStrategy, Component, EventEmitter, input, Input, Output } from '@angular/core';
import { FileData } from '../../../../../core';
import { UiTableComponent } from '../../../../../shared/components/tables/ui-table-component/ui-table-component';
import { EditButton } from '../../../../../shared/components/buttons/edit-button/edit-button';
import { DeleteButton } from '../../../../../shared/components/buttons/delete-button/delete-button';
import { DownloadButton } from '../../../../../shared/components/buttons/download-button/download-button';


@Component({
  selector: 'app-list-files-modal',
  imports: [UiTableComponent, EditButton, DeleteButton, DownloadButton],
  templateUrl: './list-files-modal.html',
  styleUrl: './list-files-modal.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListFilesModal {
  files = input.required<FileData[]>()
  isLoading = input<boolean>(false)
  canManageAccess = input<boolean>(false)

  @Output() downloadFile = new EventEmitter<FileData>()
  @Output() deleteFile = new EventEmitter<FileData>()
  @Output() editFile = new EventEmitter<FileData>()
  @Output() viewFile = new EventEmitter<FileData>()
  @Output() manageAccess = new EventEmitter<FileData>()

  onDownloadFile(file: FileData): void {
    this.downloadFile.emit(file);
  }
  onEditFile(file: FileData): void {
    this.editFile.emit(file);
  }
  onDeleteFile(file: FileData): void {
    this.deleteFile.emit(file);
  }
  onViewFile(file: FileData): void {
    this.viewFile.emit(file);
  }

  onManageAccess(file: FileData): void {
    this.manageAccess.emit(file);
  }
}
