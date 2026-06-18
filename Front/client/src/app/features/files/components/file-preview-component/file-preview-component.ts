import { ChangeDetectionStrategy, Component, computed, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { FileData } from '../../../../core';
import { ApiRoutes } from '../../../../core/routes/apiUrls';
import { SafeUrlPipe } from '../../../../shared/pipes/safe-url.pipe';
import { FileService } from '../../services';

type FileKind = 'image' | 'pdf' | 'video' | 'audio' | 'text' | 'other';

@Component({
  selector: 'app-file-preview',
  imports: [SafeUrlPipe],
  templateUrl: './file-preview-component.html',
  styleUrl: './file-preview-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FilePreviewComponent implements OnInit, OnDestroy {
  private dialogRef = inject(MatDialogRef<FilePreviewComponent>);
  private fileService = inject(FileService);
  private http = inject(HttpClient);
  protected file = inject<FileData>(MAT_DIALOG_DATA);

  objectUrl = signal<string | null>(null);
  loading = signal(true);
  loadError = signal(false);

  fileKind = computed<FileKind>(() => {
    const path = this.file.path || '';
    const ext = path.split('.').pop()?.toLowerCase() || '';
    console.log('[FilePreview] path:', path, 'ext:', ext);
    const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico'];
    const pdfExts = ['pdf'];
    const videoExts = ['mp4', 'webm', 'ogg', 'mov', 'avi', 'mkv'];
    const audioExts = ['mp3', 'wav', 'aac', 'ogg', 'flac', 'wma'];
    const textExts = ['txt', 'csv', 'json', 'xml', 'html', 'css', 'js', 'ts', 'log', 'md'];

    if (imageExts.includes(ext)) return 'image';
    if (pdfExts.includes(ext)) return 'pdf';
    if (videoExts.includes(ext)) return 'video';
    if (audioExts.includes(ext)) return 'audio';
    if (textExts.includes(ext)) return 'text';
    return 'other';
  });

  ngOnInit(): void {
    const url = ApiRoutes.files.preview(this.file.id);
    this.http.get(url, { responseType: 'blob', withCredentials: true }).subscribe({
      next: (blob) => {
        const objUrl = URL.createObjectURL(blob);
        this.objectUrl.set(objUrl);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.loadError.set(true);
      }
    });
  }

  ngOnDestroy(): void {
    const current = this.objectUrl();
    if (current) {
      URL.revokeObjectURL(current);
    }
  }

  onDownload(): void {
    this.fileService.downloadFile(this.file.id);
  }

  onClose(): void {
    const current = this.objectUrl();
    if (current) {
      URL.revokeObjectURL(current);
    }
    this.dialogRef.close();
  }
}
