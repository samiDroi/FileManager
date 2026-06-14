import { ChangeDetectionStrategy, Component, computed, inject, input, output, signal, OnInit } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { FileRequest, FileUpdateRequest, FileData, ProjectData } from '../../../../../core';
import { FileService } from '../../../services/file-service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FileTypeRequest } from '../../../../../core/domain/dtos/fileType';
import { FormInput } from '../../../../../shared/components/inputs/form-input/form-input';
import { CreateButton } from '../../../../../shared/components/buttons/create-button/create-button';
import { ConfirmButtonComponent } from '../../../../../shared/components/buttons/confirm-button-component/confirm-button-component';
import { FormComponent } from '../../../../typeFiles/components/type-files-component/form-component/form-component';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { TypeFiles } from '../../../../typeFiles/services/type-files';

@Component({
  selector: 'app-file-form',
  imports: [ReactiveFormsModule, FormInput, FormComponent, CreateButton, ConfirmButtonComponent],
  templateUrl: './file-form-component.html',
  styleUrl: './file-form-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FileFormComponent implements OnInit {
  private fb = inject(FormBuilder);
  private fileService = inject(FileService);
  private dialogRef = inject(MatDialogRef<FileFormComponent>, { optional: true });
  private dialogData = inject(MAT_DIALOG_DATA, { optional: true }) as { file?: FileData; fileTypes?: FileTypeRequest[] } | null;
  private fileTypeService = inject(TypeFiles);
  projectId = input<number>();
  fileTypes = input<FileTypeRequest[]>([]);
  projects = input<ProjectData[]>([]);
  preSelectedFile = input<File | null>(null);

  selectedFile = signal<File | null>(null);
  fileSubmit = output<FileRequest>();
  currentYear = new Date().getFullYear();
  isLoading = signal(false);
  showTypeForm = signal(false);

  isEdit = computed(() => !!this.dialogData?.file);

  fileTypesList = computed(() => this.dialogData?.fileTypes ?? this.fileTypes());

  form = this.fb.group({
    fileTitle: ['', Validators.required],
    projectId: [null as number | null],
    sameName: [false],
    year: [this.currentYear, Validators.required],
    type_id: [null as number | null, Validators.required]
  });

  constructor() {
    this.form.get('sameName')?.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe((isChecked) => {
        const titleControl = this.form.get('fileTitle');
        if (isChecked) {
          titleControl?.disable();
          titleControl?.setValue('');
        } else {
          titleControl?.enable();
        }
      });
  }

  ngOnInit(): void {
    if (this.isEdit() && this.dialogData?.file) {
      const file = this.dialogData.file;
      this.form.patchValue({
        fileTitle: file.tittle,
        projectId: file.project_id,
        year: file.year,
        type_id: file.type_id
      });
    } else if (this.preSelectedFile()) {
      this.selectedFile.set(this.preSelectedFile());
    }
  }

  onSubmitFile(): void {
    if (!this.form.valid) return;
    if (!this.isEdit() && !this.selectedFile()) return;

    const formValue = this.form.getRawValue();
    if (!formValue.fileTitle || formValue.type_id === null || formValue.year === null) return;

    this.isLoading.set(true);

    if (this.isEdit() && this.dialogData?.file) {
      const updateRequest: FileUpdateRequest = {
        tittle: formValue.fileTitle,
        year: formValue.year,
        type_id: formValue.type_id,
        project_id: formValue.projectId ?? undefined
      };
      this.fileService.editFile(this.dialogData.file.id, updateRequest);
      this.dialogRef?.close(true);
    } else if (this.selectedFile()) {
      const fileRequest: FileRequest = {
        file: this.selectedFile()!,
        project_id: formValue.projectId ?? undefined,
        tittle: formValue.fileTitle,
        year: formValue.year,
        type_id: formValue.type_id
      };
      this.fileService.uploadFile(fileRequest);
      this.dialogRef?.close(true);
    }

    this.isLoading.set(false);
  }

  onFileSelected(file: File): void {
    this.selectedFile.set(file);
    const sameName = this.form.get('sameName')?.value;
    if (sameName) {
      const nameWithoutExtension = file.name.split('.').slice(0, -1).join('.');
      this.form.patchValue({ fileTitle: nameWithoutExtension });
    }
  }

  formatFileSize(bytes?: number): string {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  }

  onCreateType(): void {
    this.showTypeForm.set(true);
  }

  handleSavedTypes(types: { name: string }[]): void {
    if (types.length > 0) {
      this.fileTypeService.createFileTypes(types);
    }
    this.showTypeForm.set(false);
  }
}
