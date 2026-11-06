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
import { EmptyStateComponent, WizardModal } from '../../../../shared/components';
import { WizardModalComponent } from './wizard-modal-component/wizard-modal-component';

import { projectService } from '../../../project/services';
import { TypeFiles } from '../../../typeFiles/services/type-files';
import { FileTypeRequest } from '../../../../core/domain/dtos/fileType';
import { ClientService } from '../../../clients/services';
import { Router } from '@angular/router';
import { AppRoutes } from '../../../../core/routes/appRoutes';
import Swal from 'sweetalert2';
import { PermissionService } from '../../../../core/services/permission-service';
import { FileAccessDialog } from '../file-access-dialog/file-access-dialog';
enum View {
  DRAG,
  FORM
}
@Component({
  selector: 'app-files-component',
  imports: [ListFilesModal, CreateButton, UiSearchBarComponent, EmptyStateComponent],
  templateUrl: './files-component.html',
  styleUrl: './files-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FilesComponent implements OnInit {
  private filesService = inject(FileService);
  private projectService = inject(projectService);
  private fileTypeService = inject(TypeFiles);
  private clientService = inject(ClientService);
  private router = inject(Router);
  private permissionService = inject(PermissionService);
  modalRef = inject(MatDialog);
  protected readonly canCreateFiles = computed(() => this.permissionService.can('crear', 'archivos'));
  protected readonly canManageAccess = computed(() => this.permissionService.can('otorgar', 'accesos'));
  isLoading = computed(() => this.filesService.state().loading);
  protected clientsLoading = computed(() => this.clientService.state().loading);
  protected hasClients = computed(() => this.clientService.state().clients.size > 0);
  protected projectsLoading = computed(() => this.projectService.state().loading);
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

  onEmptyStateClick(): void {
    this.router.navigate(['/', AppRoutes.features.dashboard, AppRoutes.features.projects]);
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
    Swal.fire({
      title: '¿Eliminar archivo?',
      text: `Se eliminará "${file.tittle}" permanentemente`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar',
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
    }).then((result) => {
      if (result.isConfirmed) {
        this.filesService.deleteFile(file);
      }
    });
  }
  onManageAccess(file: FileData): void {
    this.modalRef.open(FileAccessDialog, {
      width: '640px',
      maxWidth: '95vw',
      data: { file },
    });
  }

  onSearch(query: string): void {
    this.filesService.onSearch(query);
  }
}
