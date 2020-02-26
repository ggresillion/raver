import {Directive, ElementRef, HostListener} from '@angular/core';
import {SoundService} from '../../sound/sound.service';

@Directive({
  selector: '[appDropper]'
})
export class DropperDirective {

  constructor(
    private el: ElementRef,
    private soundService: SoundService,
  ) {
  }

  @HostListener('dragenter')
  public onDragEnter() {
    this.el.nativeElement.classList.add('drag-over');
  }

  @HostListener('dragleave')
  public onDragLeave() {
    this.el.nativeElement.classList.remove('drag-over');
  }

  @HostListener('dragover', ['$event'])
  public onDragOver(e) {
    e.preventDefault();
    e.stopPropagation();
  }

  @HostListener('drop', ['$event'])
  public onDrop(e) {
    this.el.nativeElement.classList.remove('drag-over');
    const categoryId = parseInt(e.target.id, 10);
    const soundId = parseInt(e.dataTransfer.getData('text/plain'), 10);
    this.soundService.changeCategory(soundId, categoryId)
      .subscribe(() => this.soundService.refreshSounds());
  }
}
