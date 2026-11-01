import { ChangeDetectionStrategy, Component, computed, DestroyRef, inject, OnInit, signal } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { TopbarTitleService } from '../../../../core/services';
import { NgTemplateOutlet } from '@angular/common'; // <-- Import corregido
import { FileService } from '../../services';
import { FileRequest, FileData, ProjectData } from '../../../../core';
import { CreateButton } from '../../../../shared/components/buttons/create-button/create-button';
// Importa tus componentes (rutas omitidas por brevedad)
import { ListFilesModal } from './list-files-modal/list-files-modal';
import { FileFormComponent } from '../files-component/file-form-component/file-form-component';
import { FilePreviewComponent } from '../file-preview-component/file-preview-component';
import { DragDropZoneComponent } from '../drag-drop-zone/drag-drop-zone';
import { FormComponent } from '../../../typeFiles/components/type-files-component/form-component/form-component';
import { UiSearchBarComponent } from '../../../../shared/components/search-bar/ui-search-bar-component';
import { MatDialog } from '@angular/material/dialog';
import { TypeFiles } from '../../../typeFiles/services/type-files';
import Swal from 'sweetalert2';


export enum DynamicContent {
  LIST,
  UPLOAD,
  FORM,
  TYPE_FORM
}

@Component({
  selector: 'app-project-files-modal',
  standalone: true, // Asumo que usas standalone en Angular 21
  imports: [ListFilesModal, FileFormComponent, DragDropZoneComponent, FormComponent,UiSearchBarComponent, CreateButton],
  templateUrl: './project-files-modal.html',
  styleUrl: './project-files-modal.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ProjectFilesModalComponent implements OnInit {
  private fileService = inject(FileService);
  private modalRef = inject(MatDialogRef<ProjectFilesModalComponent>);
  private editDialog = inject(MatDialog);
  private topbarTitleService = inject(TopbarTitleService);
  private destroyRef = inject(DestroyRef);
  protected data = inject(MAT_DIALOG_DATA);
  private fileTypeService = inject(TypeFiles);
  // 1. Exponemos el enum al template para usarlo en el @switch
  readonly View = DynamicContent;

  // 2. Convertimos la vista actual a un Signal reactivo
  currentView = signal<DynamicContent>(DynamicContent.LIST);
  projects = signal<ProjectData[]>([]);
  projectId = computed(() => this.data?.project?.id);
  selectedFile = signal<File | null>(null);
  files = computed(() => this.fileService.state().filesByProject.get(this.projectId()) ?? []);
  fileTypes = computed(() => this.fileTypeService.state().fileTypes);
  loading = computed(() => this.fileService.state().loading);
  
  ngOnInit(): void {
    const project = this.data?.project;
    const projectId = project?.id;

    if (project) {
      this.topbarTitleService.setProjectBreadcrumb({
        clientName: project.client_name ?? '',
        projectName: project.name ?? '',
      });
    }

    if (projectId) {
      this.fileService.getByProjectId(projectId);
      this.fileTypeService.getFileTypes();
    }
  }

  constructor() {
    this.destroyRef.onDestroy(() => this.topbarTitleService.clearProjectBreadcrumb());
  }

  // --- MÉTODOS DE NAVEGACIÓN Y ESTADO ---

  onChangeContent(): void {
    this.currentView.set(DynamicContent.UPLOAD);
    this.selectedFile.set(null);
  }

  onFileDragged(file: File): void {
    this.selectedFile.set(file);
    // Cambiamos la vista reactivamente
    this.currentView.set(DynamicContent.FORM); 
  }

  onUploadFile(file: FileRequest) {
    this.fileService.uploadFile(file);
    console.log('Archivo subido:', file);
    this.currentView.set(DynamicContent.LIST); // Volvemos a la lista después de subir
  }

  // --- OTROS MÉTODOS ---

  onEditFile(file: FileData): void {
    this.editDialog.open(FileFormComponent, {
      width: '520px',
      data: {
        file,
        fileTypes: this.fileTypes()
      }
    });
  }

  onViewFile(file: FileData): void {
    this.editDialog.open(FilePreviewComponent, {
      data: file,
      width: '95vw',
      height: '95vh',
      panelClass: 'file-preview-dialog'
    });
  }

  onClose(): void {
    this.modalRef.close();
  }

 
  onDeleteFile(file: FileData): void {
    // Cerrar el modal para que sweetalert no quede detrás
    this.modalRef.close();

    // Mostrar sweetalert con confirmación
    Swal.fire({
      title: '¿Eliminar archivo?',
      text: `¿Estás seguro de que quieres eliminar "${file.tittle}"?`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Eliminar',
      cancelButtonText: 'Cancelar',
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6'
    }).then((result) => {
      if (result.isConfirmed) {
        // Usuario confirmó: eliminar el archivo
        this.fileService.deleteFile(file);
        Swal.fire('Eliminado', 'El archivo fue eliminado correctamente', 'success');
      } else {
        // Usuario canceló: reabrir el modal con los datos originales
        this.editDialog.open(ProjectFilesModalComponent, {
          data: this.data,
           width: '850px',
      maxWidth: '95vw',
        });
      }
    });
  }

  onDownloadFile(file: FileData): void {
    this.fileService.downloadFile(file.id);
  }
  OnTypeForm(): void {
    this.currentView.set(DynamicContent.TYPE_FORM);
  }
  onSaveTypes(types: {name: string}[]): void {
    
  }
  onSearch(query: string): void {
    this.fileService.onSearch(query);
  }
}