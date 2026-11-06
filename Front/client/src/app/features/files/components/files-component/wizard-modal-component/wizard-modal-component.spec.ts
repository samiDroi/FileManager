import { ComponentFixture, TestBed } from '@angular/core/testing';
import { signal } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { FileData, FileRequest } from '../../../../../core';
import { PermissionService } from '../../../../../core/services/permission-service';
import { FileAccessDialog } from '../../file-access-dialog/file-access-dialog';
import { FileService } from '../../../services/file-service';
import { TypeFiles } from '../../../../typeFiles/services/type-files';

import { WizardModalComponent } from './wizard-modal-component';

describe('WizardModalComponent', () => {
  let component: WizardModalComponent;
  let fixture: ComponentFixture<WizardModalComponent>;
  let uploadFile: ReturnType<typeof vi.fn>;
  let openDialog: ReturnType<typeof vi.fn>;
  let canManageAccess: ReturnType<typeof vi.fn>;
  let createdFile: FileData;

  beforeEach(async () => {
    createdFile = {
      id: 42,
      project_id: 1,
      tittle: 'archivo.pdf',
      path: '/tmp/archivo.pdf',
      year: 2026,
      type_id: 1,
      file_type: 'pdf',
      created_at: '2026-11-04',
      project_name: 'Proyecto',
      can_edit: true,
      can_delete: true,
      can_download: true,
    };
    uploadFile = vi.fn((_request: FileRequest, onUploaded?: (file: FileData) => void) => {
      onUploaded?.(createdFile);
    });
    openDialog = vi.fn();
    canManageAccess = vi.fn(() => true);

    await TestBed.configureTestingModule({
      imports: [WizardModalComponent],
      providers: [
        { provide: MatDialogRef, useValue: { close: vi.fn() } },
        { provide: MatDialog, useValue: { open: openDialog } },
        { provide: FileService, useValue: { uploadFile } },
        { provide: PermissionService, useValue: { can: canManageAccess } },
        {
          provide: TypeFiles,
          useValue: {
            state: signal({ fileTypes: [] }),
            getFileTypes: vi.fn(),
            createFileTypes: vi.fn(),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(WizardModalComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('opens the access dialog after uploading when the user can grant access', () => {
    component.uploadFile({ file: new File(['x'], 'archivo.pdf'), tittle: 'archivo.pdf', year: 2026, type_id: 1, project_id: 1 });

    expect(openDialog).toHaveBeenCalledWith(FileAccessDialog, expect.objectContaining({
      data: { file: createdFile, afterUpload: true },
    }));
  });

  it('keeps the previous flow when the user cannot grant access', () => {
    canManageAccess.mockReturnValue(false);

    component.uploadFile({ file: new File(['x'], 'archivo.pdf'), tittle: 'archivo.pdf', year: 2026, type_id: 1, project_id: 1 });

    expect(openDialog).not.toHaveBeenCalled();
  });
});
