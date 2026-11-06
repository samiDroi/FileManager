import { ChangeDetectionStrategy, Component, computed, inject, signal, ChangeDetectorRef } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { FileFormComponent } from '../file-form-component/file-form-component';
import { DragDropZoneComponent } from '../../drag-drop-zone/drag-drop-zone';
import { FileTypeRequest } from '../../../../../core/domain/dtos/fileType';
import { FileRequest, ProjectData } from '../../../../../core';
import { FileService } from '../../../services/file-service';
import { TypeFiles } from '../../../../typeFiles/services/type-files';
import { PermissionService } from '../../../../../core/services/permission-service';
import { FileAccessDialog } from '../../file-access-dialog/file-access-dialog';


enum View {
  DRAG,
  FORM
}
@Component({
  selector: 'app-wizard-modal',
  imports: [FileFormComponent,DragDropZoneComponent],
  templateUrl: './wizard-modal-component.html',
  styleUrl: './wizard-modal-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class WizardModalComponent {
  private dialogRef = inject(MatDialogRef<WizardModalComponent>);
  private dialog = inject(MatDialog);
  private dialogData = inject(MAT_DIALOG_DATA, { optional: true }) as { fileTypes: FileTypeRequest[], projects: Map<number, ProjectData> } | undefined;
  private cdr = inject(ChangeDetectorRef);
  private fileService = inject(FileService);
  private fileTypeService = inject(TypeFiles);
  private permissionService = inject(PermissionService);
  protected readonly canManageAccess = computed(() => this.permissionService.can('otorgar', 'accesos'));
  readonly View = View; // Exponemos el enum al template
  actualView = signal<View>(View.DRAG);
  selectedFile = signal<File | null>(null);
  changedView = computed(() => this.actualView() === View.DRAG ? 'Drag and Drop' : 'File Form');
  fileTypes = signal<FileTypeRequest[]>(this.dialogData?.fileTypes || []);
  projects = signal(this.dialogData?.projects ? Array.from(this.dialogData.projects.values()) : []);
  onFileDragged(file: File): void {
    console.log('Archivo seleccionado:', file.name);
    this.selectedFile.set(file);
    // Cambiamos la vista reactivamente
    this.actualView.set(View.FORM);
    this.cdr.markForCheck();
  }
  onClose(): void {
    this.dialogRef.close();
  }

  uploadFile(fileRequest: FileRequest): void {
    this.fileService.uploadFile(fileRequest, createdFile => {
      const canManageAccess = this.canManageAccess();
      this.dialogRef.close();

      if (canManageAccess) {
        this.dialog.open(FileAccessDialog, {
          width: '640px',
          maxWidth: '95vw',
          data: { file: createdFile, afterUpload: true },
        });
      }
    });
  }
}
