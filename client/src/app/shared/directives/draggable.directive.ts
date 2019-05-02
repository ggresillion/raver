import {Directive, ElementRef, HostListener} from '@angular/core';

@Directive({
  selector: '[appDraggable]'
})
export class DraggableDirective {

  constructor(private el: ElementRef) {
    this.el.nativeElement.draggable = true;
  }

  @HostListener('dragstart', ['$event'])
  public onDrag(e) {
    e.dataTransfer.setData('text/plain', e.target.id);
  }
}
