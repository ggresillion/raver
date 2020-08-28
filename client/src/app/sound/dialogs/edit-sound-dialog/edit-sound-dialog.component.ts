import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Sound } from '../../../models/sound';
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { environment } from '../../../../environments/environment';
import { SoundService } from '../../sound.service';

@Component({
  selector: 'app-edit-sound-dialog',
  templateUrl: './edit-sound-dialog.component.html',
  styleUrls: ['./edit-sound-dialog.component.scss']
})
export class EditSoundDialogComponent implements OnInit {

  public readonly environment = environment;

  public form: FormGroup = new FormGroup({
    name: new FormControl('', Validators.required)
  });

  public imagePreview: string;
  private image: File;

  constructor(
    private readonly dialogRef: MatDialogRef<EditSoundDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public readonly sound: Sound,
    private readonly soundService: SoundService) { }

  public ngOnInit(): void {
    this.form.get('name').setValue(this.sound.name);
  }

  public onImageChange(event): void {
    this.image = event.target.files[0];
    const reader = new FileReader();
    reader.readAsDataURL(this.image);
    reader.onload = (_event) => {
      this.imagePreview = <string>reader.result;
    }
  }

  public save(): void {
    this.soundService.editSound(this.sound.id, this.form.value.name, this.image)
    .subscribe(() => this.dialogRef.close());
  }

  public delete(): void {
    this.soundService.deleteSound(this.sound.id)
    .subscribe(() => this.dialogRef.close());
  }

}
