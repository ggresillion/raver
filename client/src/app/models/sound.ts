import { Category } from './category';

export class Sound {
  public id: number;
  public name: string;
  public path: string;
  public category: Category;
  public image: {
    id: number,
    uuid: string
  }
}
