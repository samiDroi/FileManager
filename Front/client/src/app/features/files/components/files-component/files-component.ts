import { ChangeDetectionStrategy, Component, computed, inject, OnInit, signal } from '@angular/core';
import { ListFilesModal } from '../project-files-modal/list-files-modal/list-files-modal';
import { FileData } from '../../../../core';
import { FileService } from '../../services';
import { CreateButton } from '../../../../shared/components/buttons/create-button/create-button';
import { MAT_DIALOG_DATA, MatDialog } from '@angular/material/dialog';
import { UiSearchBarComponent } from '../../../../shared/components/search-bar/ui-search-bar-component';
import { FileFormComponent } from './file-form-component/file-form-component';
import { FilePreviewComponent } from '../file-preview-component/file-preview-component';
import { DragDropZoneComponent } from '../drag-drop-zone/drag-drop-zone';
import { ProjectFilesModalComponent } from '../project-files-modal';
import { WizardModal } from '../../../../shared/components';
import { WizardModalComponent } from './wizard-modal-component/wizard-modal-component';

import { projectService } from '../../../project/services';
import { TypeFiles } from '../../../typeFiles/services/type-files';
import { FileTypeRequest } from '../../../../core/domain/dtos/fileType';
enum View {
  DRAG,
  FORM
}
@Component({
  selector: 'app-files-component',
  imports: [ListFilesModal,CreateButton,UiSearchBarComponent],
  templateUrl: './files-component.html',
  styleUrl: './files-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FilesComponent implements OnInit {
  private filesService = inject(FileService);
  private projectService = inject(projectService);
  private fileTypeService = inject(TypeFiles);
  modalRef = inject(MatDialog);
  isLoading = signal(false);
  allFiles = this.filesService.allFiles;
  actualView = signal<View>(View.DRAG);
  selectedFile = signal<File | null>(null);
  // fileTypes = computed(() => this.fileTypeService.state().fileTypes);
  
  fileTypes = computed(() => this.fileTypeService.state().fileTypes);
  ngOnInit(): void{
    console.log(this.fileTypes());
    // Cargar todos los archivos al iniciar
    this.filesService.onSearch('');
    this.projectService.getAll();
    this.fileTypeService.getFileTypes();
  }

  onCreateFile(): void {
    const projects = this.projectService.state().projects;
    //open modal for file upload
    this.modalRef.open(WizardModalComponent, {
      width: '70%',
      maxHeight: '90vh',
      data: {
        fileTypes: this.fileTypes(),
        projects: projects
      }
    });
  }
  // constructor() {
  //   this.loadAllFiles();
  // }

  // loadAllFiles(): void {
  //   this.isLoading.set(true);
  //   this.filesService.getAll();
  //   // Los archivos se actualizarán automáticamente via el computed
  //   this.isLoading.set(false);
  // }

  onEditFile(file: FileData): void {
    this.modalRef.open(FileFormComponent, {
      width: '520px',
      data: {
        file,
        fileTypes: this.fileTypes()
      }
    });
  }

  onDownloadFile(file: FileData): void {
    this.filesService.downloadFile(file.id);
  }

  onViewFile(file: FileData): void {
    this.modalRef.open(FilePreviewComponent, {
      data: file,
      width: '95vw',
      height: '95vh',
      panelClass: 'file-preview-dialog'
    });
  }

  onDeleteFile(file: FileData): void {
    this.filesService.deleteFile(file);
  }
  onSearch(query: string): void {
    this.filesService.onSearch(query);
  }
}
