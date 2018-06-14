import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'removeExtension'
})
export class RemoveExtensionPipe implements PipeTransform {

  transform(value: string, args?: any): any {
    return value.substring(0, value.lastIndexOf('.'));
  }

}
